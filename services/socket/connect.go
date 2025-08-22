// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
)

type Connect struct {
	Base         *base.Base
	Config       Config
	Emitter      Emitter
	Broadcasters *Broadcasters
	Graceful     *base.Graceful
}

func (p *Connect) NewTxConnect(clientData *ClientSocketData, ctx context.Context) *TxConnect {
	return &TxConnect{
		BaseTx:       p.Base.NewTx(ctx),
		Config:       p.Config,
		Emitter:      p.Emitter,
		Broadcasters: p.Broadcasters,
		Graceful:     p.Graceful,
		ClientData:   clientData,
	}
}

func (p *Connect) Run(clientData *ClientSocketData, ctx context.Context) error {
	return p.Base.RunTx(p.NewTxConnect(clientData, ctx))
}

type TxConnect struct {
	*base.BaseTx
	Config       Config
	Emitter      Emitter
	ClientData   *ClientSocketData
	Broadcasters *Broadcasters
	Graceful     *base.Graceful
}

func (t *TxConnect) LoadData() (bool, error) {
	t.Graceful.WG.Add(1)

	t.ClientData.eps = NewConnectionEPS(t.Config.MaxEPS, t.Clock)
	t.ClientData.Ctx = base.GetLogger(t.Ctx).
		With().
		Str("playerID", t.ClientData.PlayerID.String()).
		Str("userID", t.ClientData.User.GetID().String()).
		WithContext(t.ClientData.Ctx)
	rooms := make([]types.SocketRoom, len(t.ClientData.User.GetRooms()))
	for i, r := range t.ClientData.User.GetRooms() {
		rooms[i] = types.SocketRoom{
			ID:      r.ID,
			Profile: &r,
		}
	}

	player := &types.Player{
		ID:     t.ClientData.PlayerID,
		RoomID: nil,
		Rooms:  rooms,
	}
	if t.ClientData.joinArgs == nil {
		err := t.Emitter.Emit(t.ClientData.Client, types.SocketMessageTypeEventV1PlayerUpdate, types.MV1PlayerUpdate{
			Player: player,
		}, t.Ctx)
		if err != nil {
			return true, err
		}
	} else {
		//	tx := s.NewTxRoomJoin(c, *data.ConnectData.Room, data.Ctx)
		//	tx.emitEvent.Player = playerDTO
		//	err := s.RunTx(tx)
		//	if service.LogErrorIfNotNil(data.Ctx, err, "join room") {
		//		service.LogErrorIfNotNil(data.Ctx, c.Close(), "close connection on join")
		//		return
		//	}
	}
	t.Broadcasters.PlayerBroadcasters.Join(t.ClientData.PlayerID, t.ClientData.PlayerID, t.ClientData.Client)
	t.Broadcasters.PrivateUserBroadcasters.Join(t.ClientData.User.GetID(), t.ClientData.PlayerID, t.ClientData.Client)
	t.Broadcasters.PublicUserBroadcasters.Join(t.ClientData.User.GetID(), t.ClientData.PlayerID, t.ClientData.Client)
	for _, r := range t.ClientData.User.GetRooms() {
		t.Broadcasters.PublicRoomBroadcasters.Join(r.ID, t.ClientData.PlayerID, t.ClientData.Client)
		t.Broadcasters.PrivateRoomBroadcasters.Join(r.ID, t.ClientData.PlayerID, t.ClientData.Client)
	}
	t.ClientData.joinArgs = nil
	base.GetLogger(WithConnectionLoggerFields(t.Ctx, t.ClientData.Ctx)).Debug().
		Msg("socket connected")
	return true, nil
}
