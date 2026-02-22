package output

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"ordo/internal/app"
	"ordo/internal/domain"
)

const (
	levelInfo  = "INFO"
	levelOK    = "OK"
	levelWarn  = "WARN"
	levelError = "ERROR"
)

const (
	colorModeAuto   = "auto"
	colorModeAlways = "always"
	colorModeNever  = "never"
)

const (
	envColorDisable = "ORDO_NO_COLORS"
	envNoLevel      = "ORDO_NO_LEVEL"
)

const (
	ansiReset  = "\x1b[0m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
	ansiGreen  = "\x1b[32m"
)

var (
	outputColorMode = colorModeAuto
	outputShowLevel = true
)

type Printer struct{}

func NewPrinter() Printer {
	return Printer{}
}

func (p Printer) Handle(errWriter io.Writer, err error) error {
	if err == nil {
		return nil
	}

	_ = PrintRootError(errWriter, err)
	return &ExitError{Code: ExitCode(err)}
}

func ParseColorMode(raw string) (string, error) {
	mode := strings.ToLower(strings.TrimSpace(raw))
	switch mode {
	case colorModeAuto, colorModeAlways, colorModeNever:
		return mode, nil
	default:
		return "", fmt.Errorf("invalid value for --color: %q (want auto, always, or never)", raw)
	}
}

func SetOutputColorMode(mode string) {
	outputColorMode = mode
}

func SetOutputShowLevel(show bool) {
	outputShowLevel = show
}

func ParseNoLevelEnv() (bool, bool, error) {
	raw, ok := os.LookupEnv(envNoLevel)
	if !ok {
		return false, false, nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return false, true, fmt.Errorf("invalid value for %s: %q (want a boolean)", envNoLevel, raw)
	}
	return parsed, true, nil
}

func shouldColorize(w io.Writer, mode string) bool {
	switch mode {
	case colorModeAlways:
		return true
	case colorModeNever:
		return false
	case colorModeAuto:
		if _, disabled := os.LookupEnv(envColorDisable); disabled {
			return false
		}
		return isTTY(w)
	default:
		return false
	}
}

func isTTY(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

func levelColor(level string) string {
	switch level {
	case levelError:
		return ansiRed
	case levelWarn:
		return ansiYellow
	case levelInfo:
		return ansiCyan
	case levelOK:
		return ansiGreen
	default:
		return ""
	}
}

func formatLevel(level string, colorEnabled bool) string {
	if !colorEnabled {
		return level
	}
	color := levelColor(level)
	if color == "" {
		return level
	}
	return color + level + ansiReset
}

func writeLevelLine(w io.Writer, level string, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	if !outputShowLevel {
		_, err := fmt.Fprintf(w, "%s\n", msg)
		return err
	}
	tag := formatLevel(level, shouldColorize(w, outputColorMode))
	_, err := fmt.Fprintf(w, "[%s] %s\n", tag, msg)
	return err
}

// PrintRootError writes a top-level CLI error using the shared output format.
func PrintRootError(w io.Writer, err error) error {
	return writeLevelLine(w, levelError, "%v", err)
}

func ExitCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrInvalidTarget):
		return 2
	case errors.Is(err, app.ErrWorkspaceNotFound):
		return 3
	case errors.Is(err, app.ErrScriptNotFound), errors.Is(err, app.ErrPackageNotFound):
		return 4
	default:
		return 1
	}
}

type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("exit code %d", e.Code)
}
