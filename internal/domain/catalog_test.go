package domain

import "testing"

func TestSupportsCatalogs(t *testing.T) {
	cases := []struct {
		manager PackageManager
		want    bool
	}{
		{ManagerNPM, false},
		{ManagerPNPM, true},
		{ManagerYarn, true},
		{ManagerBun, true},
	}

	for _, tc := range cases {
		if got := SupportsCatalogs(tc.manager); got != tc.want {
			t.Fatalf("SupportsCatalogs(%s) = %v, want %v", tc.manager, got, tc.want)
		}
	}
}

func TestParseCatalogSpec(t *testing.T) {
	tests := []struct {
		in      string
		wantPkg string
		wantVer string
		wantErr bool
	}{
		{in: "react@19.1.0", wantPkg: "react", wantVer: "19.1.0"},
		{in: "@types/node@22", wantPkg: "@types/node", wantVer: "22"},
		{in: "@scope/pkg", wantPkg: "@scope/pkg", wantVer: ""},
		{in: "react", wantPkg: "react", wantVer: ""},
		{in: "", wantErr: true},
		{in: "react@", wantErr: true},
	}

	for _, tc := range tests {
		got, err := ParseCatalogSpec(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("ParseCatalogSpec(%q) error = nil, want non-nil", tc.in)
			}
			continue
		}
		if err != nil {
			t.Fatalf("ParseCatalogSpec(%q) error = %v", tc.in, err)
		}
		if got.Package != tc.wantPkg || got.Version != tc.wantVer {
			t.Fatalf("ParseCatalogSpec(%q) = %#v, want pkg=%q ver=%q", tc.in, got, tc.wantPkg, tc.wantVer)
		}
	}
}

func TestCatalogReference(t *testing.T) {
	if got := CatalogReference(""); got != "catalog:" {
		t.Fatalf("CatalogReference(default) = %q", got)
	}
	if got := CatalogReference("react19"); got != "catalog:react19" {
		t.Fatalf("CatalogReference(named) = %q", got)
	}
}

func TestValidateCatalogName(t *testing.T) {
	if err := ValidateCatalogName("react_19"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ValidateCatalogName("React19"); err == nil {
		t.Fatal("expected validation error for uppercase")
	}
}
