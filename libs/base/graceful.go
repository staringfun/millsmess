package base

import (
	"context"
	"sync"
)

type Graceful struct {
	WG     sync.WaitGroup
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewGraceful(ctx context.Context) *Graceful {
	g := &Graceful{}
	g.Ctx, g.Cancel = context.WithCancel(ctx)
	return &Graceful{}
}

func (g *Graceful) CancelWait() {
	g.Cancel()
	g.WG.Wait()
}
