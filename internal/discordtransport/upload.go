package discordtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Discord's documented Create Message request ceiling, including MIME headers.
const maxUploadRequest = 25 << 20

type uploadFile struct {
	file *os.File
	size int64
	name string
}

type uploadAttachment struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
}

func attachmentName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	for _, b := range name {
		if !(b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '-' || b == '_' || b == '.') {
			return false
		}
	}
	return true
}

// UploadMessage streams explicitly selected regular files to the authenticated
// Create Message route. A 429 can reopen the same sections without resampling paths.
func (c *Client) UploadMessage(ctx context.Context, route, token string, body []byte) ([]byte, error) {
	segments := strings.Split(route, "/")
	if len(segments) != 6 || segments[0] != "" || segments[1] != "api" || segments[2] != "v10" ||
		segments[3] != "channels" || !snowflake(segments[4]) || segments[5] != "messages" || !validToken(token) {
		return nil, errors.New("invalid Discord upload destination")
	}
	return c.upload(ctx, route, token, body, false)
}

// UploadInteraction sends a message as the initial interaction callback.
func (c *Client) UploadInteraction(ctx context.Context, route string, body []byte) ([]byte, error) {
	initial, err := validateInteractionRoute(http.MethodPost, route)
	if err != nil || !initial {
		return nil, errors.New("invalid Discord interaction upload destination")
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return c.upload(ctx, route, "", body, true)
}

func (c *Client) upload(ctx context.Context, route, token string, body []byte, initial bool) ([]byte, error) {
	var envelope struct {
		Payload map[string]json.RawMessage `json:"payload"`
		Files   []string                   `json:"files"`
	}
	if len(body) > maxMessage || json.Unmarshal(body, &envelope) != nil ||
		envelope.Payload == nil || envelope.Payload["attachments"] != nil ||
		len(envelope.Files) < 1 || len(envelope.Files) > 10 {
		return nil, errors.New("Discord upload requires a message payload and one to ten files")
	}

	message := envelope.Payload
	if initial {
		message = nil
		if string(envelope.Payload["type"]) != "4" ||
			json.Unmarshal(envelope.Payload["data"], &message) != nil || message == nil || message["attachments"] != nil {
			return nil, errors.New("Discord interaction upload requires a type 4 response with a message")
		}
	}

	files := make([]uploadFile, 0, len(envelope.Files))
	defer func() {
		for _, part := range files {
			part.file.Close()
		}
	}()
	seen := make(map[string]bool, len(envelope.Files))
	attachments := make([]uploadAttachment, 0, len(envelope.Files))
	var size int64
	for i, path := range envelope.Files {
		name := filepath.Base(path)
		if !attachmentName(name) || seen[name] {
			return nil, errors.New("Discord upload filenames must be distinct and use letters, digits, hyphens, underscores or dots")
		}
		seen[name] = true
		before, err := os.Stat(path)
		if err != nil || !before.Mode().IsRegular() {
			return nil, errors.New("Discord upload file is not an accessible regular file")
		}
		f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			return nil, errors.New("could not open Discord upload file")
		}
		actual, err := f.Stat()
		if err != nil || !actual.Mode().IsRegular() || !os.SameFile(before, actual) {
			f.Close()
			return nil, errors.New("Discord upload file changed before it was opened")
		}
		files = append(files, uploadFile{file: f, size: actual.Size(), name: name})
		attachments = append(attachments, uploadAttachment{ID: i, Filename: name})
		if actual.Size() > maxUploadRequest-size {
			return nil, errors.New("Discord upload exceeds the 25 MiB request limit")
		}
		size += actual.Size()
	}

	metadata, err := json.Marshal(attachments)
	if err != nil {
		return nil, err
	}
	message["attachments"] = metadata
	if initial {
		envelope.Payload["data"], err = json.Marshal(message)
		if err != nil {
			return nil, err
		}
	}
	payload, err := json.Marshal(envelope.Payload)
	if err != nil {
		return nil, err
	}

	var header bytes.Buffer
	writer := multipart.NewWriter(&header)
	if err := writer.WriteField("payload_json", string(payload)); err != nil {
		return nil, err
	}
	parts := make([][]byte, 0, len(files)+1)
	for i, part := range files {
		if _, err := writer.CreateFormFile(fmt.Sprintf("files[%d]", i), part.name); err != nil {
			return nil, err
		}
		parts = append(parts, bytes.Clone(header.Bytes()))
		header.Reset()
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	parts = append(parts, header.Bytes())
	for _, part := range parts {
		size += int64(len(part))
	}
	if size > maxUploadRequest {
		return nil, errors.New("Discord upload exceeds the 25 MiB request limit")
	}
	reader := func() io.Reader {
		readers := make([]io.Reader, 0, len(files)*2+1)
		for i, part := range files {
			readers = append(readers, bytes.NewReader(parts[i]), io.NewSectionReader(part.file, 0, part.size))
		}
		return io.MultiReader(append(readers, bytes.NewReader(parts[len(files)]))...)
	}
	return c.requestBody(ctx, http.MethodPost, route, token, reader, size, writer.FormDataContentType(), initial)
}
