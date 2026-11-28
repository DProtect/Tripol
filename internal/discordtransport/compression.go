package discordtransport

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
)

const (
	deflateWindow = 32 << 10
	// Incompressible data can expand by a small number of bytes per block.
	maxCompressedMessage = maxMessage + 1024
)

var errGatewayCompression = errors.New("invalid or oversized Discord Gateway zlib stream")

type gatewayInflater struct {
	pending []byte
	window  []byte
	block   io.ReadCloser
	source  bytes.Reader
	started bool
}

// A zlib Z_SYNC_FLUSH ends at a byte boundary. Reuse the last 32 KiB of
// decompressed data when reading the next segment, and add a synthetic final
// empty block so flate can distinguish completion from truncation.
func (f *gatewayInflater) feed(chunk []byte) ([]byte, error) {
	if len(chunk) > maxCompressedMessage-len(f.pending) {
		return nil, errGatewayCompression
	}
	f.pending = append(f.pending, chunk...)
	if !bytes.HasSuffix(f.pending, []byte{0, 0, 255, 255}) {
		return nil, nil
	}
	start := 0
	if !f.started {
		const deflate = 8
		cmf, flags := f.pending[0], f.pending[1]
		if cmf&15 != deflate || cmf>>4 > 7 || flags&32 != 0 || (int(cmf)*256+int(flags))%31 != 0 {
			return nil, errGatewayCompression
		}
		start = 2
		f.started = true
	}
	f.pending = append(f.pending, 1, 0, 0, 255, 255)
	f.source.Reset(f.pending[start:])
	if f.block == nil {
		f.block = flate.NewReaderDict(&f.source, f.window)
	} else {
		if err := f.block.(flate.Resetter).Reset(&f.source, f.window); err != nil {
			return nil, errGatewayCompression
		}
	}
	payload, err := io.ReadAll(io.LimitReader(f.block, maxMessage+1))
	if err != nil || len(payload) == 0 || len(payload) > maxMessage || f.source.Len() != 0 {
		return nil, errGatewayCompression
	}
	f.pending = f.pending[:0]
	if len(payload) >= deflateWindow {
		f.window = append(f.window[:0], payload[len(payload)-deflateWindow:]...)
	} else {
		keep := min(len(f.window), deflateWindow-len(payload))
		copy(f.window, f.window[len(f.window)-keep:])
		f.window = append(f.window[:keep], payload...)
	}
	return payload, nil
}
