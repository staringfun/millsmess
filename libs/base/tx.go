// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import "context"

type Tx interface {
	Prepare() error
	DeferPrepare()
	LoadData() (bool, error)
	Publish() error
	Commit() error
	Emit()
	Apply()
}

type BaseTx struct {
	Clock         Clock
	MemoryStorage MemoryStorage
	Ctx           context.Context
}

func (t *BaseTx) Prepare() error {
	return nil
}

func (t *BaseTx) DeferPrepare() {
}

func (t *BaseTx) LoadData() (bool, error) {
	return false, nil
}

func (t *BaseTx) Publish() error {
	return nil
}

func (t *BaseTx) Commit() error {
	return nil
}

func (t *BaseTx) Emit() {
}

func (t *BaseTx) Apply() {
}

func (b *Base) NewTx(ctx context.Context) *BaseTx {
	return &BaseTx{
		Clock:         b.Clock,
		MemoryStorage: b.MemoryStorage,
		Ctx:           ctx,
	}
}

func (b *Base) RunTx(tx Tx) error {
	defer tx.DeferPrepare()
	err := tx.Prepare()
	if err != nil {
		return err
	}

	finished, err := tx.LoadData()
	if err != nil {
		return err
	}
	if finished {
		return nil
	}

	err = tx.Publish()
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	tx.Emit()

	tx.Apply()

	return nil
}
