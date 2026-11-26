package core

// A validation body still uses the real non-optimizing emitter and frame
// checks. Its code is never linked, and all inter-function references are
// symbolic. Keep one reusable body buffer instead of accumulating an entire
// throwaway image. Clear references to literal strings before reuse.
func(a *asm) resetValidationBody(){
 a.b=a.b[:0]
 clear(a.fixes);a.fixes=a.fixes[:0]
 clear(a.labels);clear(a.data);clear(a.references)
 clear(a.dataOrder);a.dataOrder=a.dataOrder[:0]
 a.serial=0
 // Keep target, trace and guard settings. The compiler resets scopes and
 // per-function frame state before compiling the next function.
}
