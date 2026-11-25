# QUANTUM 2D R2

`tripol.canvas` creates software RGB24 framebuffers using dynamic integer arrays (row-major, no alpha). Dimensions must be 1..2048 each, product <=4194304. For example, 320×180 is 57,600 pixels. Use `canvas.clear`, `canvas.rectangle` (clipped), `canvas.pixel`, and `canvas.get`. NOVA adds native checked horizontal span fills, `canvas.line`, `canvas.circle`, and `canvas.blend` (constant-opacity single pixel). Mutable operations accept dynamic arrays; do not mutate raw shape metadata or pass a const alias if immutability matters. `image.flip_x`, `flip_y`, `crop`, and `grayscale` produce copies. PPM export is limited to 256 pixels. These operations are available on Windows and Linux.

`tripol.game` is Windows x64 only, one window and one native staging buffer. `game.window(title,w,h)` opens/replaces a window (ASCII title, window size 64..2048 each). `game.running()` pumps Win32 messages and tests window lifetime; `game.key(vk)` reads a virtual-key value only when the window is foreground; `game.present(pixels,w,h)` validates RGB and shape then scales the framebuffer into the client area with GDI StretchDIBits. `game.close()` destroys the window and releases its staging buffer. `game.wait(ms)` uses the existing native sleep routine; no vsync. `events` exposes arrow, space, escape codes. Escape in the demo exits. Closing via X also exits the message loop. Call running and present regularly; no retained drawing on WM_PAINT, no dedicated render thread. Resize can cause temporary repaint/flicker. RGB-array sprites support clipped drawing and a transparent color key. There is no PNG decoding, texture filtering, text/font rendering, mouse events, audio, GPU acceleration, general window controls or FPS guarantee. The owner reported the previous Windows demo working. The NOVA staging buffer is 16 MiB and uses direct GDI presentation; the new path was cross-built and PE-inspected, not executed on Windows.

```tripol
import tripol.game
import tripol.canvas
p=canvas.create(320,180)
game.window("QUANTUM",960,540)
while game.running() {
    if game.key(27) {break}
    canvas.clear(p,1052688)
    canvas.rectangle(p,320,180,140,70,40,40,16777215)
    game.present(p,320,180)
    game.wait(16)
}
game.close()
```

Use `tripol new Demo --template game` for a 720p moving-sprite demo, or `examples/quantum_game.tri`. Use `tripol build examples/quantum_game.tri --target windows -o quantum.exe` from the source tree. Generated executable calls USER32/GDI32/KERNEL32 directly; no Go runtime in user code. Platform-specific errors exit 17, reported as TP2008 by `tripol run`.
