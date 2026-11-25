# QUANTUM R3 — experimental graphics backend
The existing Windows GDI framebuffer renderer is unchanged. A separate native WGL/OpenGL 1.1 path accepts an integer array of triangles: each triangle has nine elements (x, y, RGB per vertex), coordinates -1000..1000, colors 0..16777215, maximum 10000 triangles. It renders interpolated vertex colors, clears the background and swaps buffers. It retains the window's WGL context and releases it when game.close runs. GL failure returns false. Do not mix GDI and OpenGL in the same window: close/reopen the window before switching.

Windows OpenGL may be backed by hardware or Microsoft software rendering. No claim of actual hardware acceleration or measured FPS is made. This backend does not yet provide GPU textures, shaders or batched modern GPU commands. Linux returns false; no new runtime dependencies beyond Windows system DLLs.

```cmd
tripol build examples/projects/quantum_x/opengl.tri -o quantum-gl.exe
quantum-gl.exe
```

Expected: a colored triangle, or a cyan square on the software fallback window. Escape closes. Test window resize, repeat close/reopen and report frames and elapsed ns. The current demo uses a 16 ms sleep, so it measures paced presentation, not maximum throughput.

XInput 1.4 (Windows 8+) support is available via events.pad_connected, pad_buttons, pad_axis and pad_trigger. Controllers 0–3; axes 0–3 (left X/Y, right X/Y, -32768..32767), triggers 0–1 (0..255), buttons are the native XInput mask. -1 means unavailable, which overlaps a valid axis value; check connection separately. No vibration. Linux returns disconnected. Device availability and physical controller input require manual testing.
