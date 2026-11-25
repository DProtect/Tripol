# Standard library — 21 modules, 111 callable functions

Exact source and signatures: `tripol docs NAME`, `tripol libs`, or [stdlib](../stdlib). Counts include seven JSON scanner helper functions; use the documented high-level APIs for stable behavior. Modules are embedded and offline. In nano mode unused routines are eliminated; imports do not imply runtime dependencies.

| Module | Functions | Notes |
|---|---:|---|
| arrays | 7 | Integer arrays, copying and reversing |
| console | 3 | Prompts and lines |
| convert | 8 | Unit conversions |
| discord | 2 | Public gateway discovery and API hostname only |
| encoding | 4 | Hexadecimal encoding, not cryptography |
| errors | 5 | Checked input, division and array access |
| files | 5 | Text reads/writes; 65535-byte max read |
| http | 3 | Basic HTTP/1.0 GET, raw status and body |
| https | 1 | Windows WinHTTP GET, public endpoint only |
| json | 14 | Serialization, syntax validation and raw top-level value extraction |
| math | 12 | Checked integer math and Float64 helpers |
| network | 3 | Blocking TCP IPv4 text exchange |
| random | 3 | Seeded pseudo-random, not crypto |
| search | 4 | Linear and binary search |
| sort | 4 | Integer heapsort, returns copies |
| stats | 6 | Mean, population variance, deviation, etc. |
| system | 5 | Compiled target metadata |
| text | 9 | Byte offsets, prefix/suffix and substring search |
| time | 4 | Monotonic timer and sleep, not calendar dates |
| ui | 3 | Windows modal dialogs and text form |
| validate | 6 | ASCII validation |

`json.valid(s)` checks JSON grammar, nested arrays/objects, escapes and UTF-8 bytes up to 4096 bytes and nesting depth 32. `json.get_raw(s,key)` returns the first matching top-level value as raw JSON text, empty if invalid/missing; keys match the serializer's exact quoted spelling (no escaped-key normalization); key <=512 bytes. It does not deserialize arbitrary typed records or resolve JSON paths. Escaped surrogate code units are grammar-checked, not decoded. Internal helper signatures are subject to change. Serializer limits: quoted text <=512 bytes, integer arrays <=128 items.

`ui.form(title,label,button)` displays a fixed modal form with one edit field, submit and cancel. Empty return on cancellation/empty input; text uses Windows ANSI encoding (ASCII labels required), up to 4096 bytes, may truncate longer input. The existing `ui.message` and `ui.confirm` are also Windows-only. No custom layout/events or arbitrary windows. New form is PE-verified, not yet Windows-executed.

`https.get(host,path)` uses Windows WinHTTP with OS certificate verification, DNS, port 443, 5-second per-operation timeouts, no redirects, only successful 2xx responses, ASCII host/path, body <=65535 bytes without NUL. Not an HTTP general-purpose client: no custom headers, auth, POST, async, total deadline or Linux support; certificate/network errors are fatal. The new path is PE-verified, not executed on Windows in the development environment. `discord.gateway()` requests the public gateway discovery endpoint; no bot token, authenticated requests or WebSocket support.

`network.exchange` is a blocking IPv4 TCP text exchange (half-close, read until EOF, <=65535 bytes); no DNS/UDP/IPv6/TLS here. `http.get_raw` is plaintext, no chunk decoding or redirects; use trusted test services, never send secrets. File operations can overwrite arbitrary accessible paths, are not atomic and do not create directories. Windows file paths are ANSI. No sandbox. For full signatures and examples consult the source directly.
