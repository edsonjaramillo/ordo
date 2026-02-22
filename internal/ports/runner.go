package ports

import "context"

type Runner interface {
	Run(ctx context.Context, dir string, argv []string) error
}
