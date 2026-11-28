package httptransport

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func (q spec) requestBody(ctx context.Context, progress *transfer) (io.ReadCloser, string, error) {
	if len(q.Files) > 0 {
		return multipartBody(ctx, q.Form, q.Files, progress)
	}
	if len(q.Form) > 0 {
		values := url.Values{}
		for _, entry := range q.Form {
			values.Add(entry.Name, entry.Value)
		}
		return io.NopCloser(strings.NewReader(values.Encode())), "application/x-www-form-urlencoded", nil
	}
	if q.Body == "" {
		return nil, "", nil
	}
	kind := ""
	if q.Encoding == "json" {
		kind = "application/json"
	}
	return io.NopCloser(strings.NewReader(q.Body)), kind, nil
}

func multipartBody(ctx context.Context, fields []pair, uploads []upload, progress *transfer) (io.ReadCloser, string, error) {
	var files []*os.File
	var sizes []int64
	var total int64
	closeFiles := func() {
		for _, f := range files {
			f.Close()
		}
	}
	for _, item := range uploads {
		if item.Name == "" || strings.ContainsAny(item.Name, "\r\n\x00") || strings.ContainsAny(filepath.Base(item.Path), "\r\n\x00") {
			closeFiles()
			return nil, "", errors.New("invalid_multipart")
		}
		info, err := os.Stat(item.Path)
		if err != nil || !info.Mode().IsRegular() {
			closeFiles()
			return nil, "", errors.New("invalid_upload_file")
		}
		file, err := os.OpenFile(item.Path, os.O_RDONLY|syscall.O_NONBLOCK, 0)
		if err != nil {
			closeFiles()
			return nil, "", errors.New("invalid_upload_file")
		}
		files = append(files, file)
		opened, err := file.Stat()
		if err != nil || !os.SameFile(info, opened) {
			closeFiles()
			return nil, "", errors.New("invalid_upload_file")
		}
		sizes = append(sizes, opened.Size())
		total += opened.Size()

	}
	for _, f := range fields {
		if f.Name == "" || strings.ContainsAny(f.Name, "\r\n\x00") {
			closeFiles()
			return nil, "", errors.New("invalid_multipart")
		}
	}
	progress.uploadTotal.Store(total)
	reader, writer := io.Pipe()
	mime := multipart.NewWriter(writer)
	stop := context.AfterFunc(ctx, func() { reader.CloseWithError(ctx.Err()) })
	go func() {
		defer closeFiles()
		defer stop()
		err := writeMultipart(mime, fields, uploads, files, sizes, progress)
		writer.CloseWithError(err)
	}()
	return reader, mime.FormDataContentType(), nil
}

func writeMultipart(writer *multipart.Writer, fields []pair, uploads []upload, files []*os.File, sizes []int64, progress *transfer) error {
	for _, entry := range fields {
		if err := writer.WriteField(entry.Name, entry.Value); err != nil {
			return err
		}
	}
	for i, item := range uploads {
		dst, err := writer.CreateFormFile(item.Name, filepath.Base(item.Path))
		if err != nil {
			return err
		}
		if _, err = io.CopyN(uploadWriter{dst, progress}, files[i], sizes[i]); err != nil {
			return err
		}
	}
	return writer.Close()
}

type uploadWriter struct {
	writer   io.Writer
	progress *transfer
}

func (w uploadWriter) Write(data []byte) (int, error) {
	n, err := w.writer.Write(data)
	w.progress.uploaded.Add(int64(n))
	return n, err
}

type progressWriter struct {
	file     *os.File
	progress *transfer
}

func (w progressWriter) Write(data []byte) (int, error) {
	n, err := w.file.Write(data)
	w.progress.bytes.Add(int64(n))
	return n, err
}

func download(path string, body io.Reader, progress *transfer) (int64, string) {
	if _, err := os.Lstat(path); err == nil {
		return 0, "file_exists"
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, "file_error"
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".tripol-download-")
	if err != nil {
		return 0, "file_error"
	}
	defer os.Remove(file.Name())
	n, err := io.Copy(progressWriter{file, progress}, body)
	if err != nil {
		file.Close()
		return n, "transfer_failed"
	}
	err = file.Sync()
	closeErr := file.Close()
	if err != nil || closeErr != nil {
		return n, "file_error"
	}
	if err = os.Link(file.Name(), path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return n, "file_exists"
		}
		return n, "file_error"
	}
	return n, ""
}
