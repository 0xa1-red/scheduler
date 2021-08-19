// Package config contains the global configuration
// and the means to set and retrieve values
package config

import (
	"os"

	"hq.0xa1.red/axdx/config"
)

type Config struct {
	Test string
}

func Load(path string) (*config.Configuration, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	cfg, err := config.Load(&Config{}, fp)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
