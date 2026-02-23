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
	availability      ports.PackageManagerAvailability
}

func NewGlobalCompletionService(
	installCompletion InstallCompletionService,
	lister ports.GlobalPackageLister,
	availability ports.PackageManagerAvailability,
) GlobalCompletionService {
	return GlobalCompletionService{
		installCompletion: installCompletion,
		lister:            lister,
		availability:      availability,
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

func (s GlobalCompletionService) AvailablePackageManagers(ctx context.Context, prefix string) ([]string, error) {
	if s.availability == nil {
		return filterPrefixAndSort(domain.SupportedPackageManagers(), prefix), nil
	}

	items, err := s.availability.AvailablePackageManagers(ctx)
	if err != nil || len(items) == 0 {
		return filterPrefixAndSort(domain.SupportedPackageManagers(), prefix), nil
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
