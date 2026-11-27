package discordtransport

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	OpStart uint32 = iota + 1
	OpPoll
	OpClose
	OpREST
)

// Serve is the private, process-local protocol between a launcher and one
// generated program. Command and reply ends are inherited, never published.
func Serve(ctx context.Context, client *Client, input io.Reader, output io.Writer, diagnostics io.Writer) error {
	var gateway *Gateway
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
		case OpStart:
			if lengths[0] == 0 || lengths[0] > 4096 || lengths[1] != 0 || lengths[2] != 0 {
				return errors.New("invalid Discord bridge start")
			}
		case OpPoll, OpClose:
			if lengths != [3]int{} {
				return errors.New("invalid Discord bridge command")
			}
		case OpREST:
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
		case OpStart:
			if gateway != nil {
				err = errors.New("Discord Gateway already running")
			} else {
				gateway, err = client.StartGateway(ctx, args[0])
				reply = []byte("1")
			}
		case OpPoll:
			if gateway == nil {
				err = errors.New("Discord Gateway is closed")
			} else {
				var done bool
				reply, done, err = gateway.Poll(ctx, 100*time.Millisecond)
				if done && err == nil {
					reply = []byte("!")
				}
			}
		case OpClose:
			if gateway != nil {
				gateway.Close()
				gateway = nil
			}
			reply = []byte("1")
		case OpREST:
			method := http.MethodGet
			if len(args[2]) != 0 {
				method = http.MethodPost
			}
			reply, err = client.Request(ctx, method, string(args[0]), string(args[1]), args[2])
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
		if _, err := output.Write(header[:]); err != nil {
			return err
		}
		if len(reply) > 0 {
			if _, err := output.Write(reply); err != nil {
				return err
			}
		}
	}
}
