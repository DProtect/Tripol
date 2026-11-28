package httptransport

import (
	"compress/gzip"
	"errors"
	"io"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

func decodedBody(encoding string, body io.Reader) (io.Reader, func(), error) {
	switch strings.ToLower(strings.TrimSpace(encoding)) {
	case "", "identity":
		return body, func() {}, nil
	case "gzip":
		reader, err := gzip.NewReader(body)
		if err != nil {
			return nil, nil, err
		}
		return reader, func() { reader.Close() }, nil
	case "br":
		return brotli.NewReader(body), func() {}, nil
	case "zstd":
		reader, err := zstd.NewReader(body, zstd.WithDecoderConcurrency(1), zstd.WithDecoderMaxMemory(64<<20))
		if err != nil {
			return nil, nil, err
		}
		return reader, reader.Close, nil
	default:
		return nil, nil, errors.New("unsupported content encoding")
	}
}
