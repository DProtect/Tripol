package core

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Sentinel diagnostics carry source provenance independently of CLI formatting.
type diagnostic struct {
	Code       string `json:"code"`
	Phase      string `json:"phase"`
	Message    string `json:"message"`
	File       string `json:"file,omitempty"`
	Line       int    `json:"line,omitempty"`
	Column     int    `json:"column,omitempty"`
	Hint       string `json:"hint,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
	FixBefore  string `json:"fix_before,omitempty"`
	FixAfter   string `json:"fix_after,omitempty"`
	Detail     string `json:"detail,omitempty"`
	ExitStatus int    `json:"exit_status,omitempty"`
}

func (d diagnostic) Error() string { return d.Message }

var diagnosticJSON bool
var diagnosticTrace bool
var diagnosticSpanish bool
var entryFile string
var lineRE = regexp.MustCompile(`^line ([0-9]+): (.*)$`)

func diagnose(phase, msg, file string, line, col int) diagnostic {
	d := diagnostic{Code: "TP1000", Phase: phase, Message: msg, File: file, Line: line, Column: col}
	switch phase {
	case "lexer":
		d.Code = "TP1001"
		d.Hint = "Check the character or close the quoted string."
	case "parser":
		d.Code = "TP1002"
		d.Hint = "Check braces, parentheses, commas and statement separators."
	case "runtime":
		d.Code = "TP2000"
		d.Hint = "Check the program's error output and inputs."
	case "io":
		d.Code = "TP1006"
		d.Hint = "Check the file path, permissions and working directory."
	case "internal":
		d.Code = "TP9000"
		d.Hint = "Report this compiler failure with a minimal source example."
	default:
		switch {
		case strings.Contains(msg, "unknown variable"), strings.Contains(msg, "unknown function"), strings.Contains(msg, "unknown field"), strings.Contains(msg, "unknown method"):
			d.Code = "TP1003"
			d.Hint = "Check spelling, declaration order and imports."
		case strings.Contains(msg, "expected") && strings.Contains(msg, "got"), strings.Contains(msg, "unknown type"), strings.Contains(msg, "argument count"), strings.Contains(msg, "parameter type"):
			d.Code = "TP1004"
			d.Hint = "Check argument types, return type and function signature."
		case strings.Contains(msg, "import"), strings.Contains(msg, "module"), strings.Contains(msg, "Vault"), strings.Contains(msg, "library"):
			d.Code = "TP1005"
			d.Hint = "Check import names, paths and tripol verify."
		case strings.Contains(msg, "constant"):
			d.Code = "TP1007"
			d.Hint = "Use a mutable variable or keep the constant unchanged."
		case strings.Contains(msg, "return"):
			d.Code = "TP1008"
			d.Hint = "Ensure every required function path returns a value."
		default:
			d.Hint = "Check the source and use tripol check to validate it."
		}
	}
	return d
}
func normalizeDiagnostic(r any) diagnostic {
	if d, ok := r.(diagnostic); ok {
		return d
	}
	msg := fmt.Sprint(r)
	line := 0
	if m := lineRE.FindStringSubmatch(msg); m != nil {
		line, _ = strconv.Atoi(m[1])
		msg = m[2]
	}
	phase := "compiler"
	if _, ok := r.(*os.PathError); ok {
		phase = "io"
	}
	if strings.Contains(msg, "no such file") || strings.Contains(msg, "cannot find") || strings.Contains(msg, "permission denied") || strings.Contains(msg, "Access is denied") {
		phase = "io"
	}
	if strings.HasPrefix(msg, "internal") || strings.HasPrefix(msg, "runtime error:") {
		phase = "internal"
	}
	return diagnose(phase, msg, entryFile, line, 0)
}
func attachSource(file string, action func()) {
	defer func() {
		if r := recover(); r != nil {
			d := normalizeDiagnostic(r)
			if d.File == "" || d.File == entryFile {
				d.File = file
			}
			panic(d)
		}
	}()
	action()
}
func parseNamed(source, file string) (im []string, fs []fn, ss []stmt) {
	attachSource(file, func() { im, fs, ss = parseUnit(source) })
	for i := range fs {
		fs[i].source = file
		tagSource(fs[i].body, file)
	}
	tagSource(ss, file)
	return
}
func tagSource(ss []stmt, file string) {
	for i := range ss {
		ss[i].source = file
		tagSource(ss[i].yes, file)
		tagSource(ss[i].no, file)
		tagSource(ss[i].final, file)
	}
}
func printDiagnostic(d diagnostic) {
	if diagnosticSpanish {
		switch d.Code {
		case "TP1001":
			d.Hint = "Revisa el caracter o cierra las comillas del texto."
		case "TP1002":
			d.Hint = "Revisa llaves, parentesis, comas y separadores."
		case "TP1003":
			d.Hint = "Revisa el nombre, la declaracion y las importaciones."
		case "TP1004":
			d.Hint = "Revisa los tipos y la firma de la funcion."
		case "TP1005":
			d.Hint = "Revisa las importaciones, las rutas y tripol verify."
		case "TP1006":
			d.Hint = "Comprueba la ruta, los permisos y el directorio actual."
		case "TP1007":
			d.Hint = "Usa una variable mutable o no modifiques la constante."
		case "TP1008":
			d.Hint = "Comprueba que todos los caminos devuelvan un valor."
		case "TP1010":
			d.Hint = "Anota el tipo o usa un solo tipo para cada parametro."
		case "TP1011":
			d.Hint = "Compara los tipos de los argumentos con las firmas disponibles del metodo."
		case "TP1012":
			d.Hint = "Captura un literal escalar const o pasa el valor como argumento."
		case "TP1013":
			d.Hint = "El callback debe tener la cantidad de parametros indicada en la firma esperada."
		case "TP2001":
			d.Hint = "Comprueba la division por cero y el desbordamiento entero."
		case "TP2002":
			d.Hint = "Comprueba el indice y la longitud del arreglo."
		case "TP2003":
			d.Hint = "Revisa la condicion de assert."
		case "TP2004":
			d.Hint = "Revisa la entrada, longitud del texto y limite de memoria."
		case "TP2005":
			d.Hint = "Comprueba la ruta, permisos y longitud del archivo."
		case "TP2006":
			d.Hint = "Comprueba el servidor, puerto y limite de respuesta."
		case "TP2007":
			d.Hint = "Comprueba el rango del decimal."
		case "TP2008":
			d.Hint = "Los dialogos requieren una sesion grafica de Windows."
		case "TP9000":
			d.Hint = "Reporta el fallo con un ejemplo minimo."
		default:
			if d.Phase == "runtime" {
				d.Hint = "Revisa los datos y el mensaje de ejecucion."
			} else {
				d.Hint = "Revisa el codigo con tripol check."
			}
		}
	}
	if diagnosticJSON {
		b, _ := json.Marshal(d)
		fmt.Fprintln(os.Stderr, string(b))
		return
	}
	fmt.Fprintf(os.Stderr, "TRIPOL error [%s] (%s): %s\n", d.Code, d.Phase, d.Message)
	if d.File != "" {
		fmt.Fprintf(os.Stderr, "  --> %s", d.File)
		if d.Line > 0 {
			fmt.Fprintf(os.Stderr, ":%d", d.Line)
			if d.Column > 0 {
				fmt.Fprintf(os.Stderr, ":%d", d.Column)
			}
		}
		fmt.Fprintln(os.Stderr)
	}
	if d.Line > 0 && d.File != "" && !strings.HasPrefix(d.File, "tripol.") {
		if b, e := os.ReadFile(d.File); e == nil && len(b) <= 4<<20 {
			lines := strings.Split(string(b), "\n")
			if d.Line <= len(lines) {
				s := strings.TrimRight(lines[d.Line-1], "\r")
				if len(s) <= 240 {
					fmt.Fprintf(os.Stderr, "  %d | %s\n", d.Line, s)
					if d.Column > 0 && d.Column <= 241 {
						fmt.Fprintf(os.Stderr, "    | %s^\n", strings.Repeat(" ", d.Column-1))
					}
				}
			}
		}
	}
	if d.Hint != "" {
		label := "  Hint:"; if diagnosticSpanish {label="  Ayuda:"}; fmt.Fprintln(os.Stderr,label,d.Hint)
	}
	if d.Suggestion != "" {
		label := "  Did you mean:"; if diagnosticSpanish {label="  Quisiste decir:"}; fmt.Fprintln(os.Stderr,label,d.Suggestion+"?")
	}
	if d.FixAfter != "" {
		fmt.Fprintln(os.Stderr, "  Suggested change (not applied):")
		fmt.Fprintln(os.Stderr, "  -", d.FixBefore)
		fmt.Fprintln(os.Stderr, "  +", d.FixAfter)
	}
}
func diagnosticOptions() {
	args := os.Args[:1]
	for _, s := range os.Args[1:] {
		switch s {
		case "--diagnostics=json":
			diagnosticJSON = true
		case "--trace":
			diagnosticTrace = true
		case "--lang=es":
			diagnosticSpanish = true
		default:
			args = append(args, s)
		}
	}
	os.Args = args
}
func runtimeDiagnostic(status int) diagnostic {
	msg := fmt.Sprintf("program exited with status %d", status)
	d := diagnose("runtime", msg, entryFile, 0, 0)
	d.ExitStatus = status
	switch status {
	case 10:
		d.Code = "TP2001"
		d.Hint = "Check integer overflow and division by zero."
	case 11:
		d.Code = "TP2002"
		d.Hint = "Check the index against len(array)."
	case 12:
		d.Code = "TP2003"
		d.Hint = "An assertion failed; check its condition."
	case 13:
		d.Code = "TP2004"
		d.Hint = "Check integer input, text length and the memory limit."
	case 14:
		d.Code = "TP2005"
		d.Hint = "Check file paths, permissions and text size."
	case 15:
		d.Code = "TP2006"
		d.Hint = "Check the address, port, server and response limit."
	case 16:
		d.Code = "TP2007"
		d.Hint = "Check the floating point range."
	case 17:
		d.Code = "TP2008"
		d.Hint = "Dialogs require Windows and a working graphical session."
	}
	if status >= 100 && status <= 255 {
		d.Code = "TP2009"
		d.Hint = "An explicit exception was not caught."
	}
	return d
}

func explainDiagnostic(code string) {
	codes := map[string]string{
		"TP1000": "Compiler rejected the source or command.", "TP1001": "Unexpected source character or invalid quoted text.", "TP1002": "Invalid syntax or missing delimiter.", "TP1003": "Unknown variable, function or field.", "TP1004": "Type or argument mismatch.", "TP1005": "Import, module or package failure.", "TP1006": "File path or permission failure.", "TP1007": "Attempt to modify a constant.", "TP1008": "Invalid or missing return.", "TP2000": "Unhandled program failure.", "TP2001": "Invalid integer arithmetic (native exit 10).", "TP2002": "Index out of bounds (native exit 11).", "TP2003": "Assertion failure (native exit 12).", "TP2004": "Input, text or memory failure (native exit 13).", "TP2005": "File operation failure (native exit 14).", "TP2006": "Network operation failure (native exit 15).", "TP2007": "Float outside printable range (native exit 16).", "TP2008": "Graphical operation unavailable or failed (native exit 17).", "TP2009": "Unhandled explicit exception (native exit 100–255).", "TP9000": "Unexpected internal compiler failure."}
	codes["TP1010"] = "An inferred type is missing or conflicting."
	codes["TP1011"] = "A method overload is invalid, ambiguous or has no match."
	codes["TP1012"] = "An anonymous callback tried to capture a value without a safe constant lifetime."
	codes["TP1013"] = "A callback does not have the required parameter count."
	if s, ok := codes[strings.ToUpper(code)]; ok {
		fmt.Printf("%s: %s\n", strings.ToUpper(code), s)
		return
	}
	die("unknown diagnostic code %s", code)
}
