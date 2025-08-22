package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
)

type LeaveTx struct {
	Args   *types.MV1RoomLeave
	Reason *types.LeaveReason
	*PlayerUpdateTx
}

type Leave struct {
	Broadcasters *Broadcasters
	Emitter      Emitter
	Base         *base.Base
}

func (l *Leave) NewLeaveTx(clientData *ClientSocketData, args *types.MV1RoomLeave, ctx context.Context) base.Tx {
	return &LeaveTx{
		Args: args,
		PlayerUpdateTx: &PlayerUpdateTx{
			Broadcasters: l.Broadcasters,
			Emitter:      l.Emitter,
			UserRoomTx: &UserRoomTx{
				ClientSocketData: clientData,
				Tx: &Tx{
					BaseTx: l.Base.NewTx(ctx),
				},
			},
		},
	}
}

func (l *LeaveTx) LoadData() (bool, error) {
	finished, err := l.PlayerUpdateTx.LoadData()
	if err != nil {
		return true, err
	}
	if finished {
		return true, nil
	}
	if l.UserRooms == nil {
		return true, nil
	}
	var r *types.UserRoom
	for _, room := range l.UserRooms {
		if room.RoomID == l.Args.RoomID {
			if room.Region == l.Config.Region {
				if l.ClientSocketData.PlayerID == room.PlayerID {
					r = &room
					break
				}
			}
		}
	}
	if r == nil {
		return true, nil
	}
	_, err = l.lock.DeleteUserRoom(r.PlayerID, l.Ctx)
	if err != nil {
		return true, err
	}
	if l.Data == nil {
		l.Data = &types.MV1PlayerUpdate{}
	}
	l.Data.Key = l.Args.Key
	l.Data.LeftRoom = &l.Args.RoomID
	l.Data.LeftRoomReason = l.Reason

	if l.Msg == nil {
		l.Msg = &types.MV1PlayerMove{}
	}
	l.Msg.PlayerID = l.ClientSocketData.PlayerID
	l.Msg.UserID = l.ClientSocketData.User.GetID()
	l.Msg.LeftRoom = &l.Args.RoomID
	l.Msg.LeftRoomReason = l.Reason

	return false, nil
}

func (l *LeaveTx) Emit() {
	l.PlayerUpdateTx.Emit()
}

func (l *LeaveTx) Apply() {

}
