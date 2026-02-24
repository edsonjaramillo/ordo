package domain

import (
	"fmt"
	"regexp"
	"strings"
)

var catalogNamePattern = regexp.MustCompile(`^[a-z0-9_-]+$`)

type CatalogSpec struct {
	Package string
	Version string
}

func SupportsCatalogs(manager PackageManager) bool {
	switch manager {
	case ManagerNPM:
		return false
	case ManagerPNPM, ManagerYarn, ManagerBun:
		return true
	default:
		return false
	}
}

func ParseCatalogSpec(raw string) (CatalogSpec, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return CatalogSpec{}, fmt.Errorf("package spec cannot be empty")
	}

	idx := strings.LastIndex(value, "@")
	if idx > 0 {
		name := strings.TrimSpace(value[:idx])
		version := strings.TrimSpace(value[idx+1:])
		if name == "" {
			return CatalogSpec{}, fmt.Errorf("invalid package spec: %q", raw)
		}
		if version == "" {
			return CatalogSpec{}, fmt.Errorf("invalid package version in spec: %q", raw)
		}
		return CatalogSpec{Package: name, Version: version}, nil
	}

	return CatalogSpec{Package: value}, nil
}

func CatalogReference(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "catalog:"
	}
	return "catalog:" + trimmed
}

func ValidateCatalogName(raw string) error {
	name := strings.TrimSpace(raw)
	if name == "" {
		return fmt.Errorf("catalog name cannot be empty")
	}
	if !catalogNamePattern.MatchString(name) {
		return fmt.Errorf("invalid catalog name %q (allowed: [a-z0-9_-]+)", raw)
	}
	return nil
}
