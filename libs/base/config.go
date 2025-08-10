// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"github.com/jinzhu/configor"
	"strings"
)

type ConfigLoader interface {
	Load(paths string, data any) error
}

type DefaultConfigLoader struct {
	*configor.Configor
}

func NewDefaultConfigLoader() ConfigLoader {
	return &DefaultConfigLoader{configor.New(&configor.Config{
		Silent: true,
	})}
}

func (d *DefaultConfigLoader) Load(paths string, data any) error {
	p := strings.Split(paths, ",")
	p = append(p, "config.yaml", "config.yml")
	return d.Configor.Load(data, p...)
}
