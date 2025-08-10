// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"context"
	"fmt"
	"os"
)

type Config struct {
	InstanceID   string       `json:"instanceID" yaml:"instanceID" env:"INSTANCE_ID"`
	LoggerConfig LoggerConfig `json:"logger" yaml:"logger"`
}

type Base struct {
	Name   string
	Config Config

	ConfigLoader  ConfigLoader
	Clock         Clock
	MemoryStorage MemoryStorage
	Graceful      *Graceful

	Ctx context.Context
}

func NewBase(name string, config Config, ctx context.Context) *Base {
	loader := NewDefaultConfigLoader()
	err := loader.Load(os.Getenv("CONFIG_PATHS"), &config)
	if config.InstanceID == "" {
		config.InstanceID = fmt.Sprintf("%s-%s", name, GenerateRandomString(8))
	}
	ctx = NewLogger(config.LoggerConfig).With().
		Str("service", name).
		Str("instanceID", config.InstanceID).
		WithContext(ctx)
	LogErrorIfNotNil(ctx, err, "load config")

	return &Base{
		Name:          name,
		Config:        config,
		ConfigLoader:  loader,
		Clock:         &DefaultClock{},
		MemoryStorage: nil,
		Graceful:      NewGraceful(ctx),
		Ctx:           ctx,
	}
}
