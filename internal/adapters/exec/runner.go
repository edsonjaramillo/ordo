package execadapter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type Runner struct{}

func NewRunner() Runner {
	return Runner{}
}

func (r Runner) Run(ctx context.Context, dir string, argv []string) error {
	if len(argv) == 0 {
		return fmt.Errorf("empty command")
	}

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
