// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import "context"

type MemoryStorageTransaction interface {
	Commit() error
	Rollback() error
	MemoryStorageOperations
}

type MemoryStorageOperations interface {
}

type MemoryStorageLock interface {
	LockValue(key string, ctx context.Context) error
	LockValueRetry(key string, ctx context.Context, strategy RetryStrategy) error
	UnlockValue(key string, ctx context.Context) error
}

type MemoryStorage interface {
	MemoryStorageOperations
	MemoryStorageLock
	StartTransaction(ctx context.Context) (MemoryStorageTransaction, error)
}
