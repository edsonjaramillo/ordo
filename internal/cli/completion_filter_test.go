package cli

import (
	"context"
	"os"
	"testing"

	"ordo/internal/app"
	"ordo/internal/cli/completion"
	"ordo/internal/cli/output"
	"ordo/internal/domain"

	"github.com/spf13/cobra"
)

func TestFilterCompletedArgs(t *testing.T) {
	t.Parallel()

	got := filterCompletedArgs(
		[]string{"react", "typescript", "eslint"},
		[]string{"npm", "react", "react", "typescript"},
		1,
	)

	want := []string{"eslint"}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d (got=%#v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestFilterCompletedArgsDoesNotFilterBeforeStartIndex(t *testing.T) {
	t.Parallel()

	got := filterCompletedArgs(
		[]string{"npm", "react"},
		[]string{"npm"},
		1,
	)

	want := []string{"npm", "react"}
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d (got=%#v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestGlobalUninstallCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	completer := newTestGlobalCompleter(
		testWorkspaceIndexer{},
		testGlobalLister{packages: []string{"eslint", "prettier", "typescript"}},
	)
	cmd := newGlobalUninstallCmd(app.GlobalUninstallUseCase{}, completer, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"npm", "typescript"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "prettier"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestGlobalUpdateCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	completer := newTestGlobalCompleter(
		testWorkspaceIndexer{},
		testGlobalLister{packages: []string{"eslint", "prettier", "typescript"}},
	)
	cmd := newGlobalUpdateCmd(app.GlobalUpdateUseCase{}, completer, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"npm", "prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestGlobalInstallCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Dependencies: map[string]struct{}{
					"eslint":     {},
					"prettier":   {},
					"typescript": {},
				},
			},
		},
	}
	completer := newTestGlobalCompleter(indexer, testGlobalLister{})
	cmd := newGlobalInstallCmd(app.GlobalInstallUseCase{}, completer, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"npm", "prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestInstallCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Dependencies: map[string]struct{}{
					"eslint":     {},
					"prettier":   {},
					"typescript": {},
				},
			},
		},
	}
	cmd := newInstallCmd(app.InstallUseCase{}, newTestTargetCompleter(indexer), output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestPresetCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	cfg := []byte(`{"presets":{"web":{"devDependencies":["eslint","prettier","typescript"]}}}`)
	presetCompleter := completion.NewPresetCompleter(app.NewPresetCompletionService(testConfigStore{payload: cfg}))
	cmd := newPresetCmd(app.PresetUseCase{}, presetCompleter, completion.TargetCompleter{}, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"web", "devDependencies", "prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestPresetCompletionPresetAndBucketAreUnchanged(t *testing.T) {
	t.Parallel()

	cfg := []byte(`{"presets":{"web":{"dependencies":["react"],"devDependencies":["prettier"]}}}`)
	presetCompleter := completion.NewPresetCompleter(app.NewPresetCompletionService(testConfigStore{payload: cfg}))
	cmd := newPresetCmd(app.PresetUseCase{}, presetCompleter, completion.TargetCompleter{}, output.NewPrinter())

	presets, dir := cmd.ValidArgsFunction(cmd, []string{}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(presets) != 1 || presets[0] != "web" {
		t.Fatalf("unexpected presets: %#v", presets)
	}

	buckets, dir := cmd.ValidArgsFunction(cmd, []string{"web"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	wantBuckets := []string{"dependencies", "devDependencies"}
	if len(buckets) != len(wantBuckets) {
		t.Fatalf("len(buckets) = %d, want %d (buckets=%#v)", len(buckets), len(wantBuckets), buckets)
	}
	for i := range wantBuckets {
		if buckets[i] != wantBuckets[i] {
			t.Fatalf("buckets[%d] = %q, want %q", i, buckets[i], wantBuckets[i])
		}
	}
}

func TestCatalogPresetsCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	cfg := []byte(`{"presets":{"web":{"devDependencies":["eslint","prettier","typescript"]}}}`)
	presetCompleter := completion.NewPresetCompleter(app.NewPresetCompletionService(testConfigStore{payload: cfg}))
	cmd := newCatalogPresetsCmd(app.CatalogUseCase{}, completion.CatalogCompleter{}, presetCompleter, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"web", "devDependencies", "prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func TestCatalogPresetsCompletionPresetAndBucketAreUnchanged(t *testing.T) {
	t.Parallel()

	cfg := []byte(`{"presets":{"web":{"dependencies":["react"],"devDependencies":["prettier"]}}}`)
	presetCompleter := completion.NewPresetCompleter(app.NewPresetCompletionService(testConfigStore{payload: cfg}))
	cmd := newCatalogPresetsCmd(app.CatalogUseCase{}, completion.CatalogCompleter{}, presetCompleter, output.NewPrinter())

	presets, dir := cmd.ValidArgsFunction(cmd, []string{}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(presets) != 1 || presets[0] != "web" {
		t.Fatalf("unexpected presets: %#v", presets)
	}

	buckets, dir := cmd.ValidArgsFunction(cmd, []string{"web"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	wantBuckets := []string{"dependencies", "devDependencies"}
	if len(buckets) != len(wantBuckets) {
		t.Fatalf("len(buckets) = %d, want %d (buckets=%#v)", len(buckets), len(wantBuckets), buckets)
	}
	for i := range wantBuckets {
		if buckets[i] != wantBuckets[i] {
			t.Fatalf("buckets[%d] = %q, want %q", i, buckets[i], wantBuckets[i])
		}
	}
}

func TestUpdateCompletionFiltersAlreadySelectedTarget(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Dependencies: map[string]struct{}{
					"eslint":     {},
					"prettier":   {},
					"typescript": {},
				},
			},
		},
	}
	cmd := newUpdateCmd(app.UpdateUseCase{}, newTestTargetCompleter(indexer), output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"typescript"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0 (items=%#v)", len(items), items)
	}
}

func TestUninstallCompletionStopsAfterSingleTarget(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Dependencies: map[string]struct{}{
					"eslint":     {},
					"prettier":   {},
					"typescript": {},
				},
			},
		},
	}
	cmd := newUninstallCmd(app.UninstallUseCase{}, newTestTargetCompleter(indexer), output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"typescript"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0 (items=%#v)", len(items), items)
	}
}

func TestRunCompletionStopsAfterSingleTarget(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Scripts: map[string]string{
					"build": "vite build",
					"test":  "vitest",
				},
			},
		},
	}
	cmd := newRunCmd(app.RunUseCase{}, newTestTargetCompleter(indexer), output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"build"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0 (items=%#v)", len(items), items)
	}
}

func TestCatalogImportCompletionStopsAfterSinglePackage(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: "web",
				Dependencies: map[string]struct{}{
					"eslint": {},
				},
				DependencyVersions: map[string]string{
					"eslint": "^9.0.0",
				},
			},
		},
	}
	discovery := app.NewDiscoveryService(indexer)
	installCompletion := app.NewInstallCompletionService(discovery, nil)
	catalogCompleter := completion.NewCatalogCompleter(
		app.NewCatalogCompletionService(discovery, installCompletion, testCatalogStore{}),
	)
	cmd := newCatalogImportCmd(app.CatalogUseCase{}, catalogCompleter, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"eslint"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}
	if len(items) != 0 {
		t.Fatalf("len(items) = %d, want 0 (items=%#v)", len(items), items)
	}
}

func TestCatalogAddCompletionFiltersAlreadySelectedPackage(t *testing.T) {
	t.Parallel()

	indexer := testWorkspaceIndexer{
		infos: []domain.PackageInfo{
			{
				Dir: ".",
				Dependencies: map[string]struct{}{
					"eslint":     {},
					"prettier":   {},
					"typescript": {},
				},
			},
		},
	}
	discovery := app.NewDiscoveryService(indexer)
	installCompletion := app.NewInstallCompletionService(discovery, nil)
	catalogCompleter := completion.NewCatalogCompleter(
		app.NewCatalogCompletionService(discovery, installCompletion, testCatalogStore{}),
	)
	cmd := newCatalogAddCmd(app.CatalogUseCase{}, catalogCompleter, output.NewPrinter())

	items, dir := cmd.ValidArgsFunction(cmd, []string{"prettier"}, "")
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("directive = %v, want %v", dir, cobra.ShellCompDirectiveNoFileComp)
	}

	want := []string{"eslint", "typescript"}
	if len(items) != len(want) {
		t.Fatalf("len(items) = %d, want %d (items=%#v)", len(items), len(want), items)
	}
	for i := range want {
		if items[i] != want[i] {
			t.Fatalf("items[%d] = %q, want %q", i, items[i], want[i])
		}
	}
}

func newTestTargetCompleter(indexer testWorkspaceIndexer) completion.TargetCompleter {
	discovery := app.NewDiscoveryService(indexer)
	return completion.NewTargetCompleter(
		discovery,
		app.NewInstallCompletionService(discovery, nil),
	)
}

func newTestGlobalCompleter(indexer testWorkspaceIndexer, lister testGlobalLister) completion.GlobalCompleter {
	discovery := app.NewDiscoveryService(indexer)
	installCompletion := app.NewInstallCompletionService(discovery, nil)
	globalCompletion := app.NewGlobalCompletionService(
		installCompletion,
		lister,
		testAvailability{managers: []string{"npm"}},
	)
	return completion.NewGlobalCompleter(globalCompletion)
}

type testWorkspaceIndexer struct {
	infos []domain.PackageInfo
}

func (i testWorkspaceIndexer) Discover(context.Context) ([]domain.PackageInfo, error) {
	return i.infos, nil
}

type testGlobalLister struct {
	packages []string
}

func (l testGlobalLister) ListInstalledGlobalPackages(context.Context, domain.PackageManager) ([]string, error) {
	return l.packages, nil
}

func (l testGlobalLister) ResolveGlobalStorePaths(context.Context, domain.PackageManager) ([]string, error) {
	return nil, nil
}

type testAvailability struct {
	managers []string
}

func (a testAvailability) AvailablePackageManagers(context.Context) ([]string, error) {
	return a.managers, nil
}

type testCatalogStore struct{}

func (s testCatalogStore) UpsertCatalogEntries(context.Context, domain.PackageManager, string, map[string]string, bool) error {
	return nil
}

func (s testCatalogStore) RemoveCatalogEntries(context.Context, domain.PackageManager, string, []string) error {
	return nil
}

func (s testCatalogStore) NamedCatalogs(context.Context, domain.PackageManager) ([]string, error) {
	return nil, nil
}

func (s testCatalogStore) CatalogPackageNames(context.Context, domain.PackageManager, string) ([]string, error) {
	return nil, nil
}

type testConfigStore struct {
	payload []byte
}

func (s testConfigStore) MkdirAll(string, os.FileMode) error {
	return nil
}

func (s testConfigStore) Exists(string) (bool, error) {
	return true, nil
}

func (s testConfigStore) ReadFile(string) ([]byte, error) {
	return s.payload, nil
}

func (s testConfigStore) WriteFile(string, []byte, os.FileMode) error {
	return nil
}
