# Standard library — 55 modules

Exact source and signatures: `tripol docs NAME`, `tripol libs`, or [stdlib](../stdlib). Counts include seven JSON scanner helper functions; use the documented high-level APIs for stable behavior. Modules are embedded and offline. In nano mode unused routines are eliminated; imports do not imply runtime dependencies.

| Module | Functions | Notes |
|---|---:|---|
| binary | 8 | Validated byte arrays and u16/u32 byte-order conversion |
| arrays | 7 | Integer arrays, copying and reversing |
| console | 3 | Prompts and lines |
| convert | 8 | Unit conversions |
| discord | 34 | Experimental Windows Gateway and REST bot, discovery and protocol helpers; see Discord guide |
| env | 1 | Windows environment access |
| simd | 1 | SSE2 exact int64 array equality |
| encoding | 4 | Hexadecimal encoding, not cryptography |
| errors | 5 | Checked input, division and array access |
| files | 5 | Text reads/writes; 65535-byte max read |
| http | 3 | Basic HTTP/1.0 GET, raw status and body |
| https | 1 | Windows WinHTTP GET, public endpoint only |
| json | 15 | Serialization, syntax validation and raw top-level value extraction |
| math | 12 | Checked integer math and Float64 helpers |
| network | 3 | Blocking TCP IPv4 text exchange |
| random | 3 | Seeded pseudo-random, not crypto |
| search | 4 | Linear and binary search |
| sort | 4 | Integer heapsort, returns copies |
| stats | 6 | Mean, population variance, deviation, etc. |
| system | 5 | Compiled target metadata |
| text | 9 | Byte offsets, prefix/suffix and substring search |
| time | 6 | Monotonic timer and sleep, not calendar dates |
| ui | 7 | Windows modal dialogs and text form |
| validate | 6 | ASCII validation |

`json.valid(s)` checks JSON grammar, nested arrays/objects, escapes and UTF-8 bytes up to 4096 bytes and nesting depth 32. `json.get_raw(s,key)` returns the first matching top-level value as raw JSON text, empty if invalid/missing; keys match the serializer's exact quoted spelling (no escaped-key normalization); key <=512 bytes. It does not deserialize arbitrary typed records or resolve JSON paths. Escaped surrogate code units are grammar-checked, not decoded. Internal helper signatures are subject to change. Serializer limits: quoted text <=512 bytes, integer arrays <=128 items.

`ui.form(title,label,button)` displays a fixed modal form with one edit field, submit and cancel. Empty return on cancellation/empty input; text uses Windows ANSI encoding (ASCII labels required), up to 4096 bytes, may truncate longer input. The existing `ui.message` and `ui.confirm` are also Windows-only. No custom layout/events or arbitrary windows. The project owner demonstrated the form on Windows; the new message-box variants are PE-verified and still need Windows execution testing.

`https.get(host,path)` uses Windows WinHTTP with OS certificate verification, DNS, port 443, 5-second per-operation timeouts, no redirects, only successful 2xx responses, ASCII host/path, body <=65535 bytes without NUL. Not an HTTP general-purpose client: no custom headers, auth, POST, async, total deadline or Linux support; certificate/network errors are fatal. The new path is PE-verified, not executed on Windows in the development environment. `discord.gateway()` requests the public discovery endpoint. OMEGA adds Discord-specific authenticated REST requests and a test bot, not general authenticated HTTPS or Gateway support. See [Discord guide](discord.md).

`network.exchange` is a blocking IPv4 TCP text exchange (half-close, read until EOF, <=65535 bytes); no DNS/UDP/IPv6/TLS here. `http.get_raw` is plaintext, no chunk decoding or redirects; use trusted test services, never send secrets. File operations can overwrite arbitrary accessible paths, are not atomic and do not create directories. Windows file paths are ANSI. No sandbox. For full signatures and examples consult the source directly.

## Infinity additions

| Module | Functions | Notes |
|---|---:|---|
| bits | 6 | Nonnegative int64 bit operations; indices 0..62 |
| checksum | 3 | Adler-32, sum8, sum16; not cryptography |
| collections | 6 | Append, concatenate, insert, remove, slice, equality; return copies, fixed length |
| colors | 8 | Packed RGB24, channels 0..255, no alpha |
| compression | 3 | Integer-array run-length encoding; not general file compression |
| datetime | 7 | Gregorian calendar 1..9999, ordinal and weekday Monday=1; no clock/time zone |
| geometry | 6 | Float64 shapes, distance and interpolation |
| matrix | 6 | Flat row-major integer matrices; explicit shapes, checked math |
| sequences | 4 | Integer sequences, Fibonacci, factorial and triangular numbers |

Collections return independent arrays capped at 4,194,304 elements. They are not amortized growable vectors or dictionaries. RLE uses `[value,count,...]` pairs, may expand data, limits to 32767 runs and 65535 decoded elements, and validates all counts. `checksum` is noncryptographic and unsuitable for authentication. `ui.warning`, `ui.error`, `ui.retry`, and `ui.choose` use native message boxes; choose returns Win32 IDs Yes=6, No=7, Cancel=2. No custom widgets yet; see Quantum canvas and game APIs below.

## Quantum additions

| Module | Functions | Notes |
|---|---:|---|
| animation | 6 | Interpolation, easing and frame selection |
| buffer | 8 | Growable int buffer, explicit reassignment on push; cap 4,194,303 values |
| canvas | 9 | RGB24 framebuffer, clipping, fill, lines, circles and pixel blend; 4,194,304 pixels max |
| events | 7 | Foreground key state (Windows) and common key codes |
| game | 6 | Single Win32 window, GDI presentation, pump, key polling, close, wait |
| image | 6 | Flip, crop, grayscale, small P3 PPM export (<=256 pixels) |
| memory | 4 | Virtual reservation, used and committed/enabled page counts, not RSS |
| path | 6 | Lexical separators/basename/dirname/extension/join; no normalization |
| physics | 5 | AABB, bounds, integration, gravity; not a rigid-body solver |
| vector | 7 | Float64 vector helpers and int array scaling |

For graphics limitations, see [graphics](graphics.md). `buffer.push` must be assigned back; it grows by copying, old backing storage stays allocated. `path` does not resolve symlinks or validate filesystem paths. The game window is Windows-only and PE-inspected, actual execution pending.

## Nova additions

| Module | Functions | Notes |
|---|---:|---|
| camera | 4 | Integer coordinates and camera bounds |
| dictionary | 8 | Fixed-capacity byte-string to int64 map, linear lookup, reusable slots |
| particles | 4 | Caller-owned integer position, velocity and life arrays |
| sprite | 3 | RGB array drawing with clipping and optional transparent color |
| testing | 4 | Checked expectations for scalar values |

Sprite source and destination must not alias; copy first if needed. The map is not a general heterogeneous collection. Testing failures use recoverable assertions. Camera world conversion truncates toward zero. Particles advance by integer ticks and have no autonomous emission.

`time.now_ns()` and `time.elapsed_ns(start)` expose a monotonic timer in nanoseconds, with platform-dependent precision. The origin is unspecified; only differences should be used.

ECLIPSE adds the stateful Discord REST API (see [Discord guide](discord.md)) and native array copy/search; no new grammar.

NEXUS: `search.linear` / `arrays.contains` use AVX2 for large arrays when supported by CPU and OS, scalar otherwise. `discord.run` uses experimental Gateway; `run_rest` preserves old behavior.

## stack / queue (EVOLUTION)

Both modules expose `create(limit=32)`, `capacity`, `size`, `push`, `peek`, `pop`, `clear`. They use bounded integer backing arrays. `push` returns false when full; `peek`/`pop` assert on empty. `clear` resets logical size without freeing arena storage. Do not mutate the internal headers. See `examples/projects/collections.tri`.

## bitset / profiler (OVERDRIVE)

`bitset.create(count)` allocates a compact bounded bitset (63 bits per word, at most 1,048,576 bits). `capacity`, `contains`, `add`, `remove`, `count` expose checked operations. `add` and `remove` report whether membership changed. `profiler.start()` and `profiler.elapsed(started)` measure elapsed monotonic nanoseconds, not CPU sampling.

## QUANTUM X modules

- `heap`: bounded signed integer min-heap, create/capacity/size/push/peek/pop. No allocation after creation during valid operations. Full push returns false; empty peek/pop assert.
- `utf8`: continuation, width_at, valid, count, scalar_at. Byte offsets, strict scalar validation (overlong encodings, surrogates and out-of-range code points rejected). NUL is not representable in current strings. No normalization or grapheme segmentation.
- `csv`: escape a single field, fields/valid for a single comma-separated record with quoted fields and doubled quotes. No file/stream parser yet.

See `examples/projects/quantum_x_libraries.tri` and `memory_stress.tri`.

## 0.27.1 additional modules

- `io`: `read(path,limit=65536)` and `write(path,data)`. Binary int arrays, values 0..255; limit 65536 bytes. Invalid sizes, bytes, oversized input and OS I/O failures are fatal, matching existing file I/O. Writes replace the destination after input validation; not transactional. Reads reserve storage for the requested limit and return actual length.
- `stream`: in-memory byte cursor `create`, `position`, `remaining`, `seek`, `read_byte`, `read_u16` (little endian). Cursor is mutable, checked, and allocates only at creation.
- `regex`: `valid`, `matches`. Bounded byte-oriented full-string matching supporting literal, `.` and `*`; pattern max 128 bytes, input max 4096. No groups, captures, alternation or Unicode matching.

- `atomic`: `load(data,index)`, `exchange(data,index,value)` returns old value, `compare_exchange(data,index,expected,value)` returns bool, `add(data,index,delta)` returns new value. Aligned x86-64 loads and locked RMW instructions; checked index/arithmetic. No general thread API or claim of thread-safe runtime allocation yet.

## QUANTUM X standard-library expansion

- `crypto`: `sha256_bytes(data)` (32 digest bytes), `sha256(data)` and `sha256_text(text)` (lowercase hex). Standard SHA-256, bounded input 65536 bytes, tested against Go's standard implementation. No encryption, KDF or independent security audit.
- `audio`: `wav_mono(samples,rate=22050)`, `square_wave(hz,milliseconds,rate=22050,amplitude=5000)`, `save_wav(path,samples,rate=22050)`, `playback_supported()`, `play_file(path)`. Signed 16-bit mono PCM, 8–48 kHz, at most 32746 samples. Windows synchronous WAV playback via WinMM; other targets return false for playback. Device availability is not guaranteed. No mixer or streaming audio.
- `cli`: `option_name`, `find_option`, `has`, `value`, `integer` operate on an explicitly supplied whitespace-separated option string with `--name` and `--name=value` tokens. Last occurrence wins; no quoting or shell interpretation. Does not read OS argv.
- `process`: `id()`, `platform()`, `exit(code=0)`. Current process ID and explicit termination, checked exit code 0–255. No child spawning.
- `thread`: `supported()`, `sum(data)`, `sleep(milliseconds)`. Windows starts and joins an OS thread for an ordered checked integer reduction; caller waits, callback cannot allocate or mutate the input. Linux reports unsupported for the worker. No general thread callbacks or claims of parallel speedup.
