package app

import (
	"context"
	"sort"

	"ordo/internal/ports"
)

const defaultSuggestionLimit = 20

type InstallCompletionService struct {
	discovery DiscoveryService
	suggestor ports.PackageSuggestor
}

func NewInstallCompletionService(discovery DiscoveryService, suggestor ports.PackageSuggestor) InstallCompletionService {
	return InstallCompletionService{discovery: discovery, suggestor: suggestor}
}

func (s InstallCompletionService) WorkspaceKeys(ctx context.Context, prefix string) ([]string, error) {
	snapshot, err := s.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}
	return snapshot.WorkspaceKeys(prefix), nil
}

func (s InstallCompletionService) PackageSpecs(ctx context.Context, prefix string) ([]string, error) {
	snapshot, err := s.discovery.Snapshot(ctx)
	if err != nil {
		return nil, err
	}

	local := snapshot.DependencyNames(prefix)
	if s.suggestor == nil {
		return local, nil
	}

	remote, err := s.suggestor.Suggest(ctx, prefix, defaultSuggestionLimit)
	if err != nil {
		return local, nil
	}
	return mergeSortUnique(local, remote), nil
}

func mergeSortUnique(a []string, b []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(a)+len(b))
	for _, item := range append(a, b...) {
		if item == "" {
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
