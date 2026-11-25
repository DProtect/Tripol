package core

import (
	"fmt"
	"math/bits"
	"strings"
)

// CORE lexer: one byte-oriented scan; token text aliases the source. Preserve
// the original lexer grammar, including byte columns and CRLF behavior.
func digit(b byte) bool { return b >= '0' && b <= '9' }
func ident(b byte) bool { return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_' }
func lex(s string) []token {
	capacity := len(s) / 5
	if capacity > 8192 {
		capacity = 8192
	}
	var ts []token
	if capacity > 0 {
		ts = make([]token, 0, capacity)
	}
	line, col := 1, 1
	for i := 0; i < len(s); {
		start := i
		k := ""
		b := s[i]
		i++
		switch {
		case b == ' ' || b == '\t' || b == '\r':
			for i < len(s) && (s[i] == ' ' || s[i] == '\t' || s[i] == '\r') {
				i++
			}
		case b == '#' || (b == '/' && i < len(s) && s[i] == '/'):
			for i < len(s) && s[i] != '\n' {
				i++
			}
		case b == '\n':
			k = "nl"
		case b == '"':
			closed := false
			for i < len(s) && s[i] != '\n' {
				ch := s[i]
				i++
				if ch == '"' {
					closed = true
					break
				}
				if ch == '\\' {
					if i == len(s) || s[i] == '\n' {
						break
					}
					i++
				}
			}
			if closed {
				k = "str"
			} else {
				lexError(b, line, col)
			}
		case digit(b):
			k = "int"
			for i < len(s) && digit(s[i]) {
				i++
			}
			if i+1 < len(s) && s[i] == '.' && digit(s[i+1]) {
				i += 2
				for i < len(s) && digit(s[i]) {
					i++
				}
			}
		case ident(b):
			k = "id"
			for i < len(s) && (ident(s[i]) || digit(s[i])) {
				i++
			}
		default:
			if i < len(s) {
				switch s[start : i+1] {
				case "+=", "-=", "*=", "/=", "%=", "==", "!=", "<=", ">=", "&&", "||", "->", "=>":
					i++
					k = "op"
				}
			}
			if k == "" {
				if strings.IndexByte("+*/%<>=!(){}[],.:;-", b) >= 0 {
					k = "op"
				} else {
					lexError(b, line, col)
				}
			}
		}
		if k != "" {
			ts = append(ts, token{s[start:i], k, line, col})
		}
		if k == "nl" {
			line++
			col = 1
		} else {
			col += i - start
		}
	}
	return append(ts, token{"<eof>", "eof", line, col})
}
func lexError(b byte, line, col int) {
	panic(diagnose("lexer", fmt.Sprintf("unexpected character %q", b), "", line, col))
}

// RBP memory operand with disp8 or disp32, including positive parameter slots.
func (a *asm) memOperand(op []byte, reg byte, off int) {
	a.emit(op...)
	if off >= -128 && off <= 127 {
		a.emit(0x45|(reg<<3), byte(off))
	} else {
		a.emit(0x85 | (reg << 3))
		a.u32(off)
	}
}
func (a *asm) integerMemoryOp(op string, off int) bool {
	if r, ok := a.loopRegs[off]; ok {
		return a.regArith(op, 0, r)
	}
	switch op {
	case "+":
		a.memOperand([]byte{0x48, 0x03}, 0, off)
	case "-":
		a.memOperand([]byte{0x48, 0x2b}, 0, off)
	case "*":
		a.memOperand([]byte{0x48, 0x0f, 0xaf}, 0, off)
	default:
		return false
	}
	a.j(0x80, "err")
	return true
}

// Signed division by positive powers of two, truncating toward zero (not floor).
// Bias is at most divisor-1; adding it to a negative dividend cannot overflow.
func (a *asm) divPowerOfTwo(n int64, remainder bool) bool {
	if n <= 0 || n > 1<<30 || n&(n-1) != 0 {
		return false
	}
	if n == 1 {
		if remainder {
			a.emit(0x31, 0xc0)
		}
		return true
	}
	shift := byte(bits.TrailingZeros64(uint64(n)))
	// CQO derives sign mask; positive bias fits in 32 bits.
	// Remainder = ((x+bias)&mask)-bias, including MinInt64.
	a.emit(0x48, 0x99, 0x81, 0xe2)
	a.u32(int(n - 1))
	a.emit(0x48, 0x01, 0xd0)
	if remainder {
		a.emit(0x25)
		a.u32(int(n - 1))
		a.emit(0x48, 0x29, 0xd0)
	} else {
		a.emit(0x48, 0xc1, 0xf8, shift)
	}
	return true
}

// Only called for a proven nonnegative loop index.
func (a *asm) nonnegativePowerOfTwo(n int64, rem bool) bool {
	if n <= 0 || n > 1<<30 || n&(n-1) != 0 {
		return false
	}
	if rem {
		a.emit(0x25)
		a.u32(int(n - 1))
	} else if n != 1 {
		a.emit(0x48, 0xc1, 0xe8, byte(bits.TrailingZeros64(uint64(n))))
	}
	return true
}
