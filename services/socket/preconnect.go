// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/internal-core-api"
	"github.com/staringfun/millsmess/libs/public-core-api"
	"github.com/staringfun/millsmess/libs/types"
	"io"
	"sync"
	"sync/atomic"
)

type Preconnect struct {
	Base         *base.Base
	CoreInternal *internal_core_api.InternalCoreAPI
}

func (p *Preconnect) NewTxPreconnect(args PreconnectArgs, ctx context.Context) *TxPreconnect {
	return &TxPreconnect{
		BaseTx:       p.Base.NewTx(ctx),
		Args:         args,
		baseCtx:      p.Base.Ctx,
		CoreInternal: p.CoreInternal,
	}
}

func (p *Preconnect) Run(args PreconnectArgs, ctx context.Context) (*SocketConnectionData, error) {
	tx := p.NewTxPreconnect(args, ctx)
	err := p.Base.RunTx(tx)
	if err != nil {
		return nil, err
	}
	return tx.ConnectionData, nil
}

type PreconnectArgs struct {
	Token         string
	Cookie        string
	XForwardedFor string
	UserAgent     string
	RoomID        types.RoomID
}

type MatchArgs struct {
	RoomID types.RoomID
}

type SocketConnectionData struct {
	ctx       context.Context
	cancel    context.CancelFunc
	user      types.User
	pongedAt  atomic.Value
	matchArgs *MatchArgs
	connectMu sync.Mutex
	Writer    io.Writer
}

type TxPreconnect struct {
	*base.BaseTx
	Args           PreconnectArgs
	CoreInternal   *internal_core_api.InternalCoreAPI
	ConnectionData *SocketConnectionData
	baseCtx        context.Context
}

func (t *TxPreconnect) LoadData() (bool, error) {
	user, _, err := t.CoreInternal.FetchMeHeaders(public_core_api.FetchMeArgs{
		Token:  t.Args.Token,
		Cookie: t.Args.Cookie,
	}, map[string]string{
		"X-Forwarded-For": t.Args.XForwardedFor,
		"User-Agent":      t.Args.UserAgent,
	}, t.Ctx)
	if err != nil {
		return true, err
	}

	ctx, cancel := context.WithCancel(t.baseCtx)
	connectionData := &SocketConnectionData{
		user:   user,
		ctx:    ctx,
		cancel: cancel,
	}
	connectionData.pongedAt.Store(t.Clock.UTCTime())
	if connectionData.user.GetBannedAt() == nil {
		if !t.Args.RoomID.IsEmpty() {
			connectionData.matchArgs = &MatchArgs{RoomID: t.Args.RoomID}
		}
	}

	t.ConnectionData = connectionData
	return true, nil
}
