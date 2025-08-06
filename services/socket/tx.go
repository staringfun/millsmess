package socket

import (
	"context"
	"github.com/staringfun/millsmess/libs/base"
)

type Tx struct {
	MemoryStorageTransaction base.MemoryStorageTransaction
	*base.BaseTx
}

func (s *Socket) NewTx(ctx context.Context) *Tx {
	return &Tx{
		BaseTx: s.Base.NewTx(ctx),
	}
}

func (s *Tx) Prepare() error {
	if s.MemoryStorageTransaction != nil {
		return nil
	}
	t, err := s.MemoryStorage.StartTransaction(s.Ctx)
	if err != nil {
		return err
	}
	s.MemoryStorageTransaction = t
	return err
}

func (s *Tx) DeferPrepare() {
	if s.MemoryStorageTransaction != nil {
		s.MemoryStorageTransaction.Rollback()
		s.MemoryStorageTransaction = nil
	}
}
