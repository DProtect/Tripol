package discordtransport

import (
	"bytes"
	"compress/zlib"
	"testing"
)

func TestGatewayInflaterPreservesWindowAndHandlesSplitFlush(t *testing.T) {
	var wire bytes.Buffer
	stream := zlib.NewWriter(&wire)
	var inflater gatewayInflater
	chunk := bytes.Repeat([]byte("The Gateway sent a repeated compressed payload. 0123456789."), 800)
	for _, payload := range [][]byte{[]byte(`{"op":10,"d":{"heartbeat_interval":1000}}`), chunk, chunk, chunk[:1000]} {
		if _, err := stream.Write(payload); err != nil {
			t.Fatal(err)
		}
		if err := stream.Flush(); err != nil {
			t.Fatal(err)
		}
		data := wire.Bytes()
		// Split the sync marker across WebSocket messages.
		part, err := inflater.feed(data[:len(data)-2])
		if err != nil || part != nil {
			t.Fatalf("premature decompression: %d bytes, %v", len(part), err)
		}
		part, err = inflater.feed(data[len(data)-2:])
		if err != nil || !bytes.Equal(part, payload) {
			t.Fatalf("decompressed %d of %d bytes: %v", len(part), len(payload), err)
		}
		wire.Reset()
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestGatewayInflaterRejectsInvalidAndOversizedData(t *testing.T) {
	var inflater gatewayInflater
	if _, err := inflater.feed([]byte{0, 0, 255, 255}); err == nil {
		t.Fatal("accepted invalid zlib header")
	}
	inflater = gatewayInflater{}
	var wire bytes.Buffer
	stream := zlib.NewWriter(&wire)
	if _, err := stream.Write(bytes.Repeat([]byte("x"), maxMessage+1)); err != nil {
		t.Fatal(err)
	}
	if err := stream.Flush(); err != nil {
		t.Fatal(err)
	}
	if _, err := inflater.feed(wire.Bytes()); err == nil {
		t.Fatal("accepted oversized decompressed data")
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	inflater = gatewayInflater{}
	if _, err := inflater.feed(make([]byte, maxCompressedMessage+1)); err == nil {
		t.Fatal("accepted oversized compressed data")
	}
}
