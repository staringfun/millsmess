package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
	"sync/atomic"
	"time"
)

type ClientSocket interface {
	base.ContextWriter
	base.Closable
}

type ClientSocketData struct {
	PlayerID types.PlayerID
	Client   ClientSocket
	User     types.User

	pongedAt   atomic.Value
	joinRoomID *types.RoomID

	eps *ConnectionEPS

	Ctx    context.Context
	cancel context.CancelFunc
}

type ConnectionEPS struct {
	max float64

	checkedAt atomic.Value
	rawValue  atomic.Int32

	clock base.Clock
}

func NewConnectionEPS(max float64, clock base.Clock) *ConnectionEPS {
	return &ConnectionEPS{max: max, clock: clock}
}

func (e *ConnectionEPS) IsTooBig() bool {
	e.rawValue.Add(1)
	now := e.clock.UTCTime()
	diff := now.Sub(e.checkedAt.Load().(time.Time))
	if diff < time.Second {
		return false
	}
	value := float64(e.rawValue.Load()) / diff.Seconds()
	if value >= e.max {
		return true
	}
	e.rawValue.Store(0)
	e.checkedAt.Store(now)
	return false
}
