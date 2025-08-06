// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"github.com/jinzhu/configor"
	"strings"
)

type Config interface {
	Load(paths string, data any) error
}

type DefaultConfig struct {
}

func (d *DefaultConfig) Load(paths string, data any) error {
	p := strings.Split(paths, ",")
	return configor.Load(data, p...)
}
