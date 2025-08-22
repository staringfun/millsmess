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
	Unlock() error
}

type MemoryStorageLocker interface {
	LockValue(key string, ctx context.Context) (MemoryStorageLock, error)
	LockValueRetry(key string, ctx context.Context, strategy RetryStrategy) (MemoryStorageLock, error)
}

type MemoryStorage interface {
	MemoryStorageOperations
	MemoryStorageLocker
	StartTransaction(ctx context.Context) (MemoryStorageTransaction, error)
}
