// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"sync"
)

type Broadcaster struct {
	writers *sync.Map
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		writers: &sync.Map{},
	}
}

func (b *Broadcaster) Write(data []byte, ctx context.Context) (int, error) {
	if b == nil {
		return 0, nil
	}
	var n int
	var err error
	b.writers.Range(func(key, value any) bool {
		n, err = value.(base.ContextWriter).Write(data, ctx)
		if err != nil {
			return false
		}
		return true
	})
	return n, err
}

func (b *Broadcaster) Join(id any, writer base.ContextWriter) *Broadcaster {
	if b == nil {
		return nil
	}
	b.writers.Store(id, writer)
	return b
}

func (b *Broadcaster) Leave(id any) *Broadcaster {
	if b == nil {
		return nil
	}
	b.writers.Delete(id)
	return b
}

type BroadcasterGroup struct {
	broadcasters     *sync.Map
	itemBroadcasters *sync.Map
}

func NewBroadcasterGroup() *BroadcasterGroup {
	return &BroadcasterGroup{
		broadcasters:     &sync.Map{},
		itemBroadcasters: &sync.Map{},
	}
}

func (g *BroadcasterGroup) Join(groudID, broadcasterID any, writer base.ContextWriter) *Broadcaster {
	bb, _ := g.broadcasters.LoadOrStore(groudID, NewBroadcaster())
	b := bb.(*Broadcaster)
	b.Join(broadcasterID, writer)

	mm, _ := g.itemBroadcasters.LoadOrStore(broadcasterID, &sync.Map{})
	m := mm.(*sync.Map)
	m.Store(groudID, true)

	return b
}

func (g *BroadcasterGroup) Leave(groudID, broadcasterID any) {
	b := g.GetBroadcaster(groudID)
	if b == nil {
		return
	}

	mm, ok := g.itemBroadcasters.Load(broadcasterID)
	if ok {
		m := mm.(*sync.Map)
		m.Delete(groudID)
	}

	b.writers.Delete(broadcasterID)

	isEmpty := true
	b.writers.Range(func(_, _ any) bool {
		isEmpty = false
		return false
	})
	if isEmpty {
		bb, loaded := g.broadcasters.LoadAndDelete(groudID)
		if loaded {
			bb.(*Broadcaster).writers.Range(func(_, value any) bool {
				g.Join(groudID, broadcasterID, value.(base.ContextWriter))
				return true
			})
		}
	}
}

func (g *BroadcasterGroup) LeaveAll(broadcasterID any) {
	mm, loaded := g.itemBroadcasters.LoadAndDelete(broadcasterID)
	if !loaded {
		return
	}
	m := mm.(*sync.Map)
	m.Range(func(key, _ any) bool {
		g.Leave(key, broadcasterID)
		return true
	})
}

func (g *BroadcasterGroup) GetBroadcaster(id any) *Broadcaster {
	b, ok := g.broadcasters.Load(id)
	if !ok {
		return nil
	}
	return b.(*Broadcaster)
}

type Broadcasters struct {
	UserBroadcasters *BroadcasterGroup
}

func NewBroadcasters() *Broadcasters {
	return &Broadcasters{
		UserBroadcasters: NewBroadcasterGroup(),
	}
}
