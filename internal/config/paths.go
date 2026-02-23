package config

import (
	"os"
	"path/filepath"
)

func WorkingDir() (string, error) {
	return os.Getwd()
}

func ConfigHome() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return xdg, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config"), nil
}

func OrdoConfigDir() (string, error) {
	configHome, err := ConfigHome()
	if err != nil {
		return "", err
	}
	return filepath.Join(configHome, "ordo"), nil
}

func OrdoConfigPath() (string, error) {
	configDir, err := OrdoConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "ordo.json"), nil
}
