package ports

import "os"

type ConfigStore interface {
	MkdirAll(path string, perm os.FileMode) error
	Exists(path string) (bool, error)
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
}
