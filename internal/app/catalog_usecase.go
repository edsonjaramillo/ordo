package app

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"ordo/internal/domain"
	"ordo/internal/ports"
)

type CatalogAddRequest struct {
	Packages  []string
	Workspace string
	Force     bool
}

type CatalogsAddRequest struct {
	Name      string
	Packages  []string
	Workspace string
	Force     bool
}

type CatalogRemoveRequest struct {
	Packages []string
}

type CatalogsRemoveRequest struct {
	Name     string
	Packages []string
}

type CatalogSyncRequest struct{}

type CatalogUseCase struct {
	discovery DiscoveryService
	catalogs  ports.CatalogStore
	manifests ports.ManifestStore
	versions  ports.PackageVersionResolver
}

func NewCatalogUseCase(
	discovery DiscoveryService,
	catalogs ports.CatalogStore,
	manifests ports.ManifestStore,
	versions ports.PackageVersionResolver,
) CatalogUseCase {
	return CatalogUseCase{
		discovery: discovery,
		catalogs:  catalogs,
		manifests: manifests,
		versions:  versions,
	}
}

func (u CatalogUseCase) RunAdd(ctx context.Context, req CatalogAddRequest) error {
	return u.applyAdd(ctx, "", req.Packages, req.Workspace, req.Force)
}

func (u CatalogUseCase) RunNamedAdd(ctx context.Context, req CatalogsAddRequest) error {
	if err := domain.ValidateCatalogName(req.Name); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidCatalogName, err)
	}
	return u.applyAdd(ctx, strings.TrimSpace(req.Name), req.Packages, req.Workspace, req.Force)
}

func (u CatalogUseCase) RunRemove(ctx context.Context, req CatalogRemoveRequest) error {
	return u.applyRemove(ctx, "", req.Packages)
}

func (u CatalogUseCase) RunNamedRemove(ctx context.Context, req CatalogsRemoveRequest) error {
	if err := domain.ValidateCatalogName(req.Name); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidCatalogName, err)
	}
	return u.applyRemove(ctx, strings.TrimSpace(req.Name), req.Packages)
}

func (u CatalogUseCase) RunSync(ctx context.Context, _ CatalogSyncRequest) error {
	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	if !domain.SupportsCatalogs(snapshot.Manager) {
		return fmt.Errorf("%w: %s", ErrCatalogUnsupported, snapshot.Manager)
	}

	packages, err := u.catalogs.CatalogPackageNames(ctx, snapshot.Manager, "")
	if err != nil {
		return err
	}
	if len(packages) == 0 {
		return nil
	}

	workspaces := sortedWorkspaceInfos(snapshot.ByWorkspace)
	for _, workspace := range workspaces {
		if err := u.manifests.RewriteCatalogReferencesExistingOnly(ctx, workspace.Dir, "", packages); err != nil {
			return err
		}
	}

	return nil
}

func (u CatalogUseCase) applyAdd(ctx context.Context, name string, rawPackages []string, workspace string, force bool) error {
	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	if !domain.SupportsCatalogs(snapshot.Manager) {
		return fmt.Errorf("%w: %s", ErrCatalogUnsupported, snapshot.Manager)
	}

	target, err := resolveInstallTargetPackage(snapshot, workspace)
	if err != nil {
		return err
	}

	resolved, err := u.resolvePackageVersions(ctx, rawPackages)
	if err != nil {
		return err
	}

	if err := u.catalogs.UpsertCatalogEntries(ctx, snapshot.Manager, name, resolved, force); err != nil {
		if strings.Contains(err.Error(), "catalog conflict") {
			return fmt.Errorf("%w: %v", ErrCatalogConflict, err)
		}
		return err
	}

	packages := sortedPackageNames(resolved)
	return u.manifests.RewriteCatalogReferences(ctx, target.Dir, name, packages)
}

func (u CatalogUseCase) applyRemove(ctx context.Context, name string, rawPackages []string) error {
	snapshot, err := u.discovery.Snapshot(ctx)
	if err != nil {
		return err
	}

	if !domain.SupportsCatalogs(snapshot.Manager) {
		return fmt.Errorf("%w: %s", ErrCatalogUnsupported, snapshot.Manager)
	}

	packages, err := validateRemovePackages(rawPackages)
	if err != nil {
		return err
	}

	return u.catalogs.RemoveCatalogEntries(ctx, snapshot.Manager, name, packages)
}

func (u CatalogUseCase) resolvePackageVersions(ctx context.Context, rawPackages []string) (map[string]string, error) {
	items := trimNonEmpty(rawPackages)
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one package is required")
	}

	resolved := make(map[string]string, len(items))
	for _, item := range items {
		spec, err := domain.ParseCatalogSpec(item)
		if err != nil {
			return nil, err
		}

		if spec.Version == "" {
			if u.versions == nil {
				return nil, fmt.Errorf("latest version resolver is not configured")
			}
			version, err := u.versions.LatestVersion(ctx, spec.Package)
			if err != nil {
				return nil, err
			}
			spec.Version = version
		}

		resolved[spec.Package] = spec.Version
	}
	return resolved, nil
}

func validateRemovePackages(rawPackages []string) ([]string, error) {
	items := trimNonEmpty(rawPackages)
	if len(items) == 0 {
		return nil, fmt.Errorf("at least one package is required")
	}

	out := make([]string, 0, len(items))
	seen := map[string]struct{}{}
	for _, item := range items {
		spec, err := domain.ParseCatalogSpec(item)
		if err != nil {
			return nil, err
		}
		if spec.Version != "" {
			return nil, fmt.Errorf("remove accepts package names only: %s", item)
		}
		if _, ok := seen[spec.Package]; ok {
			continue
		}
		seen[spec.Package] = struct{}{}
		out = append(out, spec.Package)
	}
	return out, nil
}

func sortedPackageNames(items map[string]string) []string {
	names := make([]string, 0, len(items))
	for name := range items {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func sortedWorkspaceInfos(items map[string]domain.PackageInfo) []domain.PackageInfo {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	ordered := make([]domain.PackageInfo, 0, len(keys))
	for _, key := range keys {
		ordered = append(ordered, items[key])
	}
	return ordered
}
