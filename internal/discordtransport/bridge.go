package discordtransport

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tripol/internal/httptransport"
)

const (
	OpStart uint32 = iota + 1
	OpPoll
	OpClose
	OpREST
	OpInteraction
	OpInteractionEdit
	OpRESTGet
	OpRESTPost
	OpRESTPut
	OpRESTPatch
	OpRESTDelete
	OpRESTUpload
	OpInteractionUpload
	OpCooldown
	OpAuthorizeModeration
	OpChannelLock
	OpWebhook
	OpWebhookUpload
	OpHTTP
)

// Serve is the private, process-local protocol between a launcher and one
// generated program. Command and reply ends are inherited, never published.
func Serve(ctx context.Context, client *Client, input io.Reader, output io.Writer, diagnostics io.Writer) error {
	input = bufio.NewReaderSize(input, 8192)
	writer := bufio.NewWriterSize(output, maxMessage+4)
	var gateway *Gateway
	var httpService httptransport.Service
	continuations := make(httpContinuations)
	var servedHTTP bool
	defer httpService.Close()
	defer continuations.close(&httpService)
	var callbackDeadline time.Time
	var cooldown cooldowns
	var channelLock channelLocks
	defer cooldown.close()
	defer func() {
		if gateway != nil {
			gateway.Close()
		}
	}()
	var hdr [16]byte
	for {
		if _, err := io.ReadFull(input, hdr[:]); err != nil {
			if err == io.EOF || ctx.Err() != nil {
				return nil
			}
			return err
		}
		op := binary.LittleEndian.Uint32(hdr[:4])
		var lengths [3]int
		for i := range lengths {
			length := binary.LittleEndian.Uint32(hdr[4+4*i:])
			if length > maxMessage {
				return errors.New("oversized Discord bridge request")
			}
			lengths[i] = int(length)
		}
		switch op {
		case OpHTTP:
			if lengths[0] == 0 || lengths[0] > 8 || lengths[1] != 0 || lengths[2] == 0 {
				return errors.New("invalid HTTP bridge request")
			}
		case OpStart:
			if lengths[0] == 0 || lengths[0] > 4096 || lengths[1] != 0 || lengths[2] != 0 {
				return errors.New("invalid Discord bridge start")
			}
		case OpPoll, OpClose:
			if lengths != [3]int{} {
				return errors.New("invalid Discord bridge command")
			}
		case OpInteraction, OpInteractionEdit, OpInteractionUpload:
			if lengths[0] == 0 || lengths[0] > 1024 || lengths[1] != 0 || (op != OpInteractionEdit && lengths[2] == 0) {
				return errors.New("invalid Discord bridge callback")
			}
		case OpWebhook, OpWebhookUpload:
			if lengths[0] == 0 || lengths[0] > 20 || lengths[1] == 0 || lengths[1] > 512 || lengths[2] == 0 {
				return errors.New("invalid Discord webhook bridge request")
			}
		case OpChannelLock:
			if lengths[0] == 0 || lengths[0] > 20 || lengths[1] == 0 || lengths[1] > 512 || lengths[2] > 6 {
				return errors.New("invalid Discord channel lock bridge request")
			}
		case OpAuthorizeModeration:
			if lengths[0] != 0 || lengths[1] == 0 || lengths[1] > 512 || lengths[2] == 0 || lengths[2] > 4096 {
				return errors.New("invalid Discord moderation bridge request")
			}
		case OpCooldown:
			if lengths[0] == 0 || lengths[0] > 9 || lengths[1] != 0 || lengths[2] == 0 || lengths[2] > 16384 {
				return errors.New("invalid Discord cooldown bridge request")
			}
		case OpREST, OpRESTGet, OpRESTPost, OpRESTPut, OpRESTPatch, OpRESTDelete, OpRESTUpload:
			if lengths[0] == 0 || lengths[0] > 1024 || lengths[1] == 0 || lengths[1] > 512 {
				return errors.New("invalid Discord bridge REST command")
			}
		default:
			return errors.New("unknown Discord bridge command")
		}
		var args [3][]byte
		for i, length := range lengths {
			if length == 0 {
				continue
			}
			args[i] = make([]byte, length)
			if _, err := io.ReadFull(input, args[i]); err != nil {
				return err
			}
		}
		var reply []byte
		var err error
		switch op {
		case OpHTTP:
			switch string(args[0]) {
			case "route":
				reply = []byte("0")
				if continuations.attach(&httpService, args[2]) {
					reply = []byte("1")
				}
			case "complete":
				reply = []byte("0")
				if continuations.complete(ctx, client, args[2]) {
					reply = []byte("1")
				}
			case "queue":
				reply = []byte("0")
				if continuations.edit(ctx, client, args[2], true) {
					reply = []byte("1")
				}
			case "cancelcb":
				id, _ := strconv.Atoi(string(args[2]))
				reply = []byte("0")
				if continuations.cancel(&httpService, id) {
					reply = []byte("1")
				}
			default:
				reply = httpService.Handle(ctx, string(args[0]), args[2])
			}
		case OpStart:
			callbackDeadline = time.Time{}
			if gateway != nil {
				err = errors.New("Discord Gateway already running")
			} else {
				gateway, err = client.StartGateway(ctx, args[0])
				reply = []byte("1")
			}
		case OpPoll:
			cooldown.prune(time.Now())
			if gateway == nil {
				err = errors.New("Discord Gateway is closed")
			} else {
				if servedHTTP {
					if event, ready := gateway.queued(); ready {
						reply = event.raw
						callbackDeadline = event.callbackDeadline
						servedHTTP = false
						break
					}
				}
				if continuation := continuations.event(&httpService); len(continuation) != 0 {
					reply = continuation
					callbackDeadline = time.Time{}
					servedHTTP = true
					break
				}
				servedHTTP = false
				var done bool
				var event gatewayEvent
				event, done, err = gateway.poll(ctx, 100*time.Millisecond)
				reply = event.raw
				if len(reply) != 0 {
					callbackDeadline = event.callbackDeadline
				}
				if done && err == nil {
					reply = []byte("!")
				}
			}
		case OpClose:
			callbackDeadline = time.Time{}
			servedHTTP = false
			continuations.close(&httpService)
			if gateway != nil {
				gateway.Close()
				gateway = nil
			}
			reply = []byte("1")
		case OpInteraction, OpInteractionEdit, OpInteractionUpload:
			requestCtx := ctx
			cancel := func() {}
			if !callbackDeadline.IsZero() && strings.HasPrefix(string(args[0]), "/api/v10/interactions/") {
				requestCtx, cancel = context.WithDeadline(ctx, callbackDeadline)
			}
			if op == OpInteractionUpload {
				reply, err = client.UploadInteraction(requestCtx, string(args[0]), args[2])
			} else if op == OpInteractionEdit && len(args[2]) == 0 {
				reply, err = client.DeleteInteractionFollowup(ctx, string(args[0]))
			} else if op == OpInteractionEdit {
				reply, err = client.EditInteraction(ctx, string(args[0]), args[2])
			} else {
				reply, err = client.Interaction(requestCtx, string(args[0]), args[2])
			}
			cancel()
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				reply, err = nil, nil
			}
		case OpREST:
			method := http.MethodGet
			if len(args[2]) != 0 {
				method = http.MethodPost
			}
			reply, err = client.Request(ctx, method, string(args[0]), string(args[1]), args[2])
		case OpRESTGet, OpRESTPost, OpRESTPut, OpRESTPatch, OpRESTDelete:
			methods := [...]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
			reply, err = client.Request(ctx, methods[op-OpRESTGet], string(args[0]), string(args[1]), args[2])
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				reply, err = nil, nil
			}
		case OpWebhook, OpWebhookUpload:
			reply, err = client.executeWebhook(ctx, string(args[0]), string(args[1]), args[2], op == OpWebhookUpload)
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				reply, err = nil, nil
			}
		case OpChannelLock:
			if lockErr := channelLock.change(ctx, client, string(args[0]), string(args[1]), string(args[2])); lockErr != nil {
				fmt.Fprintln(diagnostics, lockErr)
			} else {
				reply = []byte("1")
			}
		case OpAuthorizeModeration:
			requestCtx := ctx
			cancel := func() {}
			if !callbackDeadline.IsZero() {
				requestCtx, cancel = context.WithDeadline(ctx, callbackDeadline.Add(-250*time.Millisecond))
			}
			result, checkErr := client.authorizeModeration(requestCtx, string(args[1]), args[2])
			cancel()
			if checkErr != nil {
				fmt.Fprintln(diagnostics, checkErr)
			}
			reply = []byte(result)
		case OpCooldown:
			reply, err = cooldown.respond(string(args[0]), args[2])
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				reply, err = nil, nil
			}
		case OpRESTUpload:
			reply, err = client.UploadMessage(ctx, string(args[0]), string(args[1]), args[2])
			if err != nil {
				fmt.Fprintln(diagnostics, err)
				reply, err = nil, nil
			}
		}
		if err == nil && (len(reply) > maxMessage || bytes.IndexByte(reply, 0) >= 0) {
			err = errors.New("invalid Discord bridge reply")
		}
		var header [4]byte
		if err != nil {
			fmt.Fprintln(diagnostics, err)
			binary.LittleEndian.PutUint32(header[:], ^uint32(0))
			reply = nil
		} else {
			binary.LittleEndian.PutUint32(header[:], uint32(len(reply)))
		}
		if _, err := writer.Write(header[:]); err != nil {
			return err
		}
		if len(reply) > 0 {
			if _, err := writer.Write(reply); err != nil {
				return err
			}
		}
		if err := writer.Flush(); err != nil {
			return err
		}
	}
}
