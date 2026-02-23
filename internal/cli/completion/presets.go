package completion

import (
	"context"

	"ordo/internal/app"
)

type PresetCompleter struct {
	presets app.PresetCompletionService
}

func NewPresetCompleter(presets app.PresetCompletionService) PresetCompleter {
	return PresetCompleter{presets: presets}
}

func (c PresetCompleter) PresetNames(ctx context.Context, prefix string) ([]string, error) {
	return c.presets.PresetNames(ctx, prefix)
}

func (c PresetCompleter) Buckets(ctx context.Context, preset string, prefix string) ([]string, error) {
	return c.presets.Buckets(ctx, preset, prefix)
}

func (c PresetCompleter) BucketPackages(ctx context.Context, preset string, bucket string, prefix string) ([]string, error) {
	return c.presets.BucketPackages(ctx, preset, bucket, prefix)
}
