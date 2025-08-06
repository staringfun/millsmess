// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/internal-core-api"
)

type Socket struct {
	*base.Base
	Preconnect   *Preconnect
	CoreInternal *internal_core_api.InternalCoreAPI
}

func NewSocket() *Socket {
	b := &base.Base{
		Clock:         &base.DefaultClock{},
		MemoryStorage: nil,
		Ctx:           context.Background(),
	}
	coreInternal := &internal_core_api.InternalCoreAPI{}
	return &Socket{
		Base:         b,
		CoreInternal: coreInternal,
		Preconnect: &Preconnect{
			Base:         b,
			CoreInternal: coreInternal,
		},
	}
}
