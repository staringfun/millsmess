// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/internal-core-api"
)

type Config struct {
	MaxEPS float64 `json:"maxEPS" yaml:"maxEPS" env:"MAX_EPS" default:"30"`
	base.Config
}

type Socket struct {
	*base.Base
	CoreInternal *internal_core_api.InternalCoreAPI

	Emitter      Emitter
	Broadcasters *Broadcasters

	Preconnect *Preconnect
	Connect    *Connect
}

func NewSocket(config Config, ctx context.Context) *Socket {
	b := base.NewBase("socket", config.Config, ctx)
	coreInternal := &internal_core_api.InternalCoreAPI{}
	emitter := &DefaultEmitter{}
	return &Socket{
		Base:         b,
		CoreInternal: coreInternal,
		Emitter:      emitter,
		Broadcasters: NewBroadcasters(),
		Preconnect: &Preconnect{
			Base:         b,
			CoreInternal: coreInternal,
		},
		Connect: &Connect{
			Base:     b,
			Config:   config,
			Emitter:  emitter,
			Graceful: b.Graceful,
		},
	}
}
