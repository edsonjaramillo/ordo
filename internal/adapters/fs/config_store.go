package fs

import "os"

type ConfigStore struct{}

func NewConfigStore() ConfigStore {
	return ConfigStore{}
}

func (c ConfigStore) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (c ConfigStore) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (c ConfigStore) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (c ConfigStore) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}
