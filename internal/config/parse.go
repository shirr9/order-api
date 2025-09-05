package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"io/fs"
	"os"
)

func Load(path string) (*Config, error) {
	if err := godotenv.Load(); err != nil {
		//logging mb changed
		return nil, fmt.Errorf("no .env file found: %w", err)
	}

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("config file not found at %s: %w", path, err)
		}
		return nil, err
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
