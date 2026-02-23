package app

import (
	"context"
	"encoding/json"
	"fmt"

	"ordo/internal/config"
	"ordo/internal/domain"
	"ordo/internal/ports"
)

type InitRequest struct {
	DefaultPackageManager string
}

type InitUseCase struct {
	configStore ports.ConfigStore
}

const initConfigSchemaURL = "https://raw.githubusercontent.com/edsonjaramillo/ordo/refs/heads/main/schema.json"

type initConfigPayload struct {
	Schema                string `json:"$schema"`
	DefaultPackageManager string `json:"defaultPackageManager"`
}

func NewInitUseCase(configStore ports.ConfigStore) InitUseCase {
	return InitUseCase{configStore: configStore}
}

func (u InitUseCase) Run(_ context.Context, req InitRequest) error {
	manager, err := domain.ParsePackageManager(req.DefaultPackageManager)
	if err != nil {
		return err
	}

	configDir, err := config.OrdoConfigDir()
	if err != nil {
		return err
	}
	configPath, err := config.OrdoConfigPath()
	if err != nil {
		return err
	}

	exists, err := u.configStore.Exists(configPath)
	if err != nil {
		return err
	}
	if exists {
		return ErrConfigAlreadyExists
	}

	if err := u.configStore.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(initConfigPayload{
		Schema:                initConfigSchemaURL,
		DefaultPackageManager: string(manager),
	}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	payload = append(payload, '\n')
	return u.configStore.WriteFile(configPath, payload, 0o644)
}
