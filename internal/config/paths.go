package config

import "os"

func WorkingDir() (string, error) {
	return os.Getwd()
}
