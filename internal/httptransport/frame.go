package httptransport

import (
	"encoding/json"
	"strconv"
)

// The private pipe already preserves record lengths. Keep validated UTF-8
// response text as the final field instead of quoting and unquoting the body.
// File and header values stay JSON-escaped because they can contain line breaks.
func encodeTextFrame(out result) []byte {
	headers, _ := json.Marshal(out.Headers)
	file := []byte(`""`)
	if out.File != "" {
		file, _ = json.Marshal(out.File)
	}
	frame := make([]byte, 0, 96+len(file)+len(headers)+len(out.Text))
	frame = append(frame, '@')
	for _, value := range [...]int64{
		int64(out.Status), out.Bytes, out.Elapsed, int64(out.Attempts), out.Uploaded, out.UploadTotal,
	} {
		frame = strconv.AppendInt(frame, value, 10)
		frame = append(frame, '\n')
	}
	frame = append(frame, out.Protocol...)
	frame = append(frame, '\n')
	frame = append(frame, out.Error...)
	frame = append(frame, '\n')
	frame = append(frame, file...)
	frame = append(frame, '\n')
	frame = append(frame, headers...)
	frame = append(frame, '\n')
	frame = append(frame, out.Text...)
	if len(frame) > maxFrame {
		return encode(result{Status: out.Status, Protocol: out.Protocol, Bytes: out.Bytes,
			Elapsed: out.Elapsed, Attempts: out.Attempts, Error: "response_too_large"})
	}
	return frame
}
