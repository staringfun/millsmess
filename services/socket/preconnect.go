// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"errors"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/internal-core-api"
	"github.com/staringfun/millsmess/libs/public-core-api"
	"github.com/staringfun/millsmess/libs/types"
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

func (p *Preconnect) Run(args PreconnectArgs, ctx context.Context) (*ClientSocketData, error) {
	tx := p.NewTxPreconnect(args, ctx)
	err := p.Base.RunTx(tx)
	if err != nil {
		return nil, err
	}
	return tx.ClientData, nil
}

type PreconnectArgs struct {
	Token         string
	Cookie        string
	XForwardedFor string
	UserAgent     string
	JoinRoomID    types.RoomID
}

type TxPreconnect struct {
	*base.BaseTx
	Args         PreconnectArgs
	CoreInternal *internal_core_api.InternalCoreAPI
	ClientData   *ClientSocketData
	baseCtx      context.Context
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
		if errors.Is(err, public_core_api.UnauthorizedError) {
			return true, nil
		}
		if errors.Is(err, public_core_api.FetchMeError) {
			return true, nil
		}
		return true, err
	}

	bannedAt := user.GetBannedAt()
	if bannedAt != nil {
		if t.Clock.UTCTime().After(*bannedAt) {
			return true, nil
		}
	}

	ctx, cancel := context.WithCancel(t.baseCtx)
	playerID := base.GeneratePlayerID()
	ctx = WithPlayerID(ctx, playerID)
	ctx = WithUserID(ctx, user.GetID())
	t.ClientData = &ClientSocketData{
		PlayerID:   playerID,
		User:       user,
		Ctx:        ctx,
		cancel:     cancel,
		joinRoomID: &t.Args.JoinRoomID,
	}
	t.ClientData.pongedAt.Store(t.Clock.UTCTime())

	return true, nil
}
