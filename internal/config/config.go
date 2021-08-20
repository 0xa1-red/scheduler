// Package config contains the global configuration
// and the means to set and retrieve values
package config

import (
	"os"

	"hq.0xa1.red/axdx/config"
	"hq.0xa1.red/axdx/scheduler/internal/platform/api"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database"
	"hq.0xa1.red/axdx/scheduler/internal/platform/database/redis"
	"hq.0xa1.red/axdx/scheduler/internal/platform/nats"
)

var cfg *config.Configuration

type Config struct {
	API struct {
		Address string `yaml:"address"`
	} `yaml:"api"`
	Database struct {
		Kind  string `yaml:"kind" default:"redis"`
		Redis struct {
			Address  string `yaml:"address"`
			Password string `yaml:"password"`
			Database int    `yaml:"database"`
			Retries  int    `yaml:"retries"`
			Interval string `yaml:"interval"`
		} `yaml:"redis"`
	}
	NATS struct {
		Address        string `yaml:"address"`
		DefaultSubject string `yaml:"default_subject"`
	} `yaml:"nats"`
}

func load(path string) (*config.Configuration, error) {
	if cfg == nil {
		fp, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}

		c, err := config.Load(&Config{}, fp)
		if err != nil {
			return nil, err
		}

		cfg = c
	}

	return cfg, nil
}

func ConfigurePackages(configPath string) error {
	cfg, err := load(configPath)
	if err != nil {
		return err
	}

	if address, ok := cfg.Get("api.address").(string); ok && address != "" {
		api.SetAddress(address)
	}

	if kind, ok := cfg.Get("database.kind").(string); ok && kind != "" {
		database.SetBackend(kind)
	}
	switch database.Backend() {
	case database.KindRedis:
		if address, ok := cfg.Get("database.redis.address").(string); ok && address != "" {
			redis.SetAddress(address)
		}
		if password, ok := cfg.Get("database.redis.password").(string); ok && password != "" {
			redis.SetPassword(password)
		}
		if database, ok := cfg.Get("database.redis.database").(int); ok {
			redis.SetDatabase(database)
		}
		if retries, ok := cfg.Get("database.redis.retries").(int); ok {
			redis.SetRetries(retries)
		}
		if interval, ok := cfg.Get("database.redis.interval").(string); ok && interval != "" {
			redis.SetInterval(interval)
		}
	}

	if address, ok := cfg.Get("nats.address").(string); ok && address != "" {
		nats.SetAddress(address)
	}
	if subject, ok := cfg.Get("nats.default_subject").(string); ok && subject != "" {
		nats.SetDefaultSubject(subject)
	}

	return nil
}
