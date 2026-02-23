package domain

import (
	"fmt"
	"strings"
)

type PresetBucket string

const (
	BucketDependencies         PresetBucket = "dependencies"
	BucketDevDependencies      PresetBucket = "devDependencies"
	BucketPeerDependencies     PresetBucket = "peerDependencies"
	BucketOptionalDependencies PresetBucket = "optionalDependencies"
)

func SupportedPresetBuckets() []string {
	return []string{
		string(BucketDependencies),
		string(BucketDevDependencies),
		string(BucketPeerDependencies),
		string(BucketOptionalDependencies),
	}
}

func ParsePresetBucket(raw string) (PresetBucket, error) {
	value := strings.TrimSpace(raw)
	switch value {
	case string(BucketDependencies):
		return BucketDependencies, nil
	case string(BucketDevDependencies):
		return BucketDevDependencies, nil
	case string(BucketPeerDependencies):
		return BucketPeerDependencies, nil
	case string(BucketOptionalDependencies):
		return BucketOptionalDependencies, nil
	default:
		return "", fmt.Errorf("unsupported preset bucket: %q (supported: %s)", raw, strings.Join(SupportedPresetBuckets(), ", "))
	}
}

func BucketInstallOptions(bucket PresetBucket) InstallOptions {
	switch bucket {
	case BucketDevDependencies:
		return InstallOptions{Dev: true}
	case BucketPeerDependencies:
		return InstallOptions{Peer: true}
	case BucketOptionalDependencies:
		return InstallOptions{Optional: true}
	default:
		return InstallOptions{}
	}
}
