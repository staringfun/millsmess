package socket

import (
	"github.com/staringfun/millsmess/libs/base"
	internal_core_api "github.com/staringfun/millsmess/libs/internal-core-api"
	public_core_api "github.com/staringfun/millsmess/libs/public-core-api"
	"github.com/staringfun/millsmess/libs/types"
)

type Tx struct {
	CoreInternal             *internal_core_api.InternalCoreAPI
	MemoryStorageTransaction base.MemoryStorageTransaction
	*base.BaseTx
}

func (s *Tx) Prepare() error {
	err := s.BaseTx.Prepare()
	if err != nil {
		return err
	}
	if s.MemoryStorageTransaction != nil {
		return nil
	}
	t, err := s.MemoryStorage.StartTransaction(s.Ctx)
	if err != nil {
		return err
	}
	s.MemoryStorageTransaction = t
	return nil
}

func (s *Tx) DeferPrepare() {
	s.BaseTx.DeferPrepare()
	if s.MemoryStorageTransaction != nil {
		s.MemoryStorageTransaction.Rollback()
		s.MemoryStorageTransaction = nil
	}
}

func (s *Tx) Commit() error {
	err := s.BaseTx.Commit()
	if err != nil {
		return err
	}
	if s.MemoryStorageTransaction != nil {
		return s.MemoryStorageTransaction.Commit()
	}
	return nil
}

type UserRoomTx struct {
	ClientSocketData *ClientSocketData
	lock             public_core_api.CoreLock
	*Tx
}

func (s *UserRoomTx) Prepare() error {
	if s.lock == nil {
		l, err := s.CoreInternal.LockUserRoom(s.ClientSocketData.User.GetID(), s.Ctx)
		if err != nil {
			return err
		}
		s.lock = l
	}
	return s.Tx.Prepare()
}

func (s *UserRoomTx) DeferPrepare() {
	if s.lock != nil {
		s.lock.Unlock(s.Ctx)
		s.lock = nil
	}
	s.Tx.DeferPrepare()
}

type PlayerUpdateTx struct {
	Data         *types.MV1PlayerUpdate
	Msg          *types.MV1PlayerMove
	Broadcasters *Broadcasters
	Emitter      Emitter
	UserRooms    []types.UserRoom
	*UserRoomTx
}

func (p *PlayerUpdateTx) LoadData() (bool, error) {
	if p.UserRooms != nil {
		return false, nil
	}
	rooms, err := p.CoreInternal.FetchUserRooms(p.ClientSocketData.User.GetID(), p.Ctx)
	if err != nil {
		return true, err
	}
	p.UserRooms = rooms
	return p.UserRoomTx.LoadData()
}

func (p *PlayerUpdateTx) Publish() error {
	if p.Msg == nil {
		return nil
	}
	p.BaseTx
	p.Emitter.Emit(
		p.Broadcasters.PlayerBroadcasters.GetBroadcaster(p.ClientSocketData.PlayerID),
		types.SocketMessageTypeEventV1PlayerUpdate,
		p.Data,
		p.Ctx,
	)
}

func (p *PlayerUpdateTx) Emit() {
	if p.Data == nil {
		return
	}
	p.Emitter.Emit(
		p.Broadcasters.PlayerBroadcasters.GetBroadcaster(p.ClientSocketData.PlayerID),
		types.SocketMessageTypeEventV1PlayerUpdate,
		p.Data,
		p.Ctx,
	)
}
