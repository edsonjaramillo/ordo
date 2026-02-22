package domain

import "testing"

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name      string
		in        string
		workspace string
		target    string
		wantErr   bool
	}{
		{name: "root", in: "build", target: "build"},
		{name: "workspace", in: "ui/build", workspace: "ui", target: "build"},
		{name: "empty", in: "", wantErr: true},
		{name: "too many segments", in: "a/b/c", wantErr: true},
		{name: "missing workspace", in: "/build", wantErr: true},
		{name: "missing name", in: "ui/", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTarget(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Workspace != tc.workspace || got.Name != tc.target {
				t.Fatalf("unexpected parse result: %+v", got)
			}
		})
	}
}
