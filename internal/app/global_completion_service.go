package app

import (
	"context"
	"sort"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type GlobalCompletionService struct {
	installCompletion InstallCompletionService
	lister            ports.GlobalPackageLister
}

func NewGlobalCompletionService(
	installCompletion InstallCompletionService,
	lister ports.GlobalPackageLister,
) GlobalCompletionService {
	return GlobalCompletionService{
		installCompletion: installCompletion,
		lister:            lister,
	}
}

func (s GlobalCompletionService) InstallPackageSpecs(ctx context.Context, prefix string) ([]string, error) {
	return s.installCompletion.PackageSpecs(ctx, prefix)
}

func (s GlobalCompletionService) InstalledGlobalPackages(ctx context.Context, manager domain.PackageManager, prefix string) ([]string, error) {
	if s.lister == nil {
		return []string{}, nil
	}

	items, err := s.lister.ListInstalledGlobalPackages(ctx, manager)
	if err != nil {
		return []string{}, nil
	}

	return filterPrefixAndSort(items, prefix), nil
}

func filterPrefixAndSort(items []string, prefix string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if item == "" {
			continue
		}
		if prefix != "" && !strings.HasPrefix(item, prefix) {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}
