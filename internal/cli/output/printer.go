package output

import (
	"errors"
	"fmt"
	"io"

	"ordo/internal/app"
	"ordo/internal/domain"
)

type Printer struct {
	err io.Writer
}

func NewPrinter(err io.Writer) Printer {
	return Printer{err: err}
}

func (p Printer) Handle(err error) error {
	if err == nil {
		return nil
	}

	_, _ = fmt.Fprintf(p.err, "Error: %v\n", err)
	return &ExitError{Code: ExitCode(err)}
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
