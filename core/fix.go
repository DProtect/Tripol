package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Guard Fix: conservative edit-distance suggestion, never changes source.
func (c *compiler) suggest(msg string) string {
	var target string
	candidates := map[string]bool{}
	if strings.HasPrefix(msg, "unknown variable ") {
		target = strings.TrimPrefix(msg, "unknown variable ")
		for _, scope := range c.scopes {
			for n := range scope {
				candidates[n] = true
			}
		}
	}
	if strings.HasPrefix(msg, "unknown function ") {
		target = strings.TrimPrefix(msg, "unknown function ")
		for n := range c.fs {
			candidates[n] = true
		}
		for _, n := range []string{"input", "input_int", "parse_int", "len", "sqrt", "sum", "min", "max", "abs", "text", "read_text", "write_text", "append_text"} {
			candidates[n] = true
		}
	}
	if len(target) < 3 || len(target) > 128 {
		return ""
	}
	best, dist, tie := "", 3, false
	for n := range candidates {
		if strings.HasPrefix(n, "@") || len(n) > 128 {
			continue
		}
		d := editDistance(target, n)
		if d < dist {
			best, dist, tie = n, d, false
		} else if d == dist {
			tie = true
		}
	}
	if tie || dist > 2 {
		return ""
	}
	return best
}
func editDistance(a, b string) int {
	row := make([]int, len(b)+1)
	for j := range row {
		row[j] = j
	}
	for i := range a {
		prev := row[0]
		row[0] = i + 1
		for j := range b {
			old := row[j+1]
			cost := 0
			if a[i] != b[j] {
				cost = 1
			}
			row[j+1] = min(prev+cost, row[j]+1, old+1)
			prev = old
		}
	}
	return row[len(b)]
}

// A preview never writes to disk. Restrict to a single exact unqualified token
// on the reported line; comments and string literals are excluded by the lexer.
func replacementPreview(d diagnostic, source string) (string, string, bool) {
	if d.Line < 1 || d.Suggestion == "" {
		return "", "", false
	}
	prefix := "unknown variable "
	if strings.HasPrefix(d.Message, "unknown function ") {
		prefix = "unknown function "
	}
	if !strings.HasPrefix(d.Message, prefix) {
		return "", "", false
	}
	target := strings.TrimPrefix(d.Message, prefix)
	if strings.Contains(target, ".") || strings.Contains(d.Suggestion, ".") {
		return "", "", false
	}
	lines := strings.Split(source, "\n")
	if d.Line > len(lines) {
		return "", "", false
	}
	line := strings.TrimSuffix(lines[d.Line-1], "\r")
	position := -1
	for _, tok := range lex(line) {
		if tok.k == "id" && tok.s == target {
			if position >= 0 {
				return "", "", false
			}
			position = tok.col - 1
		}
	}
	if position < 0 {
		return "", "", false
	}
	return line, line[:position] + d.Suggestion + line[position+len(target):], true
}

func previewFix(file string) {
	defer func() {
		if r := recover(); r != nil {
			d := normalizeDiagnostic(r)
			// Only preview a real accessible source file inside the project; never edit.
			if d.File != "" && !strings.HasPrefix(d.File, "tripol.") {
				base, be := filepath.Abs(filepath.Dir(file))
				actual, ae := filepath.Abs(d.File)
				if be == nil && ae == nil {
					rel, re := filepath.Rel(base, actual)
					if re == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
						if data, e := os.ReadFile(actual); e == nil && len(data) <= 4<<20 {
							if before, after, ok := replacementPreview(d, string(data)); ok {
								d.FixBefore = before
								d.FixAfter = after
							}
						}
					}
				}
			}
			panic(d)
		}
	}()
	fs, ss, _ := loadProject(file)
	generate(fs, ss, "linux", false, nil)
	fmt.Println("No compiler errors found. Source unchanged.")
}
