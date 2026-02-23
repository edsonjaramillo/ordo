package domain

import "testing"

func TestParsePresetBucket(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    PresetBucket
		wantErr bool
	}{
		{name: "dependencies", input: "dependencies", want: BucketDependencies},
		{name: "devDependencies", input: "devDependencies", want: BucketDevDependencies},
		{name: "peerDependencies", input: "peerDependencies", want: BucketPeerDependencies},
		{name: "optionalDependencies", input: "optionalDependencies", want: BucketOptionalDependencies},
		{name: "invalid", input: "prodDependencies", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParsePresetBucket(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("ParsePresetBucket() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestBucketInstallOptions(t *testing.T) {
	tests := []struct {
		name   string
		bucket PresetBucket
		want   InstallOptions
	}{
		{name: "dependencies", bucket: BucketDependencies, want: InstallOptions{}},
		{name: "devDependencies", bucket: BucketDevDependencies, want: InstallOptions{Dev: true}},
		{name: "peerDependencies", bucket: BucketPeerDependencies, want: InstallOptions{Peer: true}},
		{name: "optionalDependencies", bucket: BucketOptionalDependencies, want: InstallOptions{Optional: true}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := BucketInstallOptions(tc.bucket)
			if got != tc.want {
				t.Fatalf("BucketInstallOptions() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
