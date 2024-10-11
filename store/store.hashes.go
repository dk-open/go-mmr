package store

import (
	"context"
	"github.com/dk-open/go-mmr/types"
	"sync"
)

type memoryHashSource[TH types.HashType] struct {
	sync.RWMutex
	indexes map[TH][]byte
}

func MemoryHashSource[TH types.HashType]() IHashSource[TH] {
	return &memoryHashSource[TH]{
		indexes: make(map[TH][]byte),
	}
}

func (a *memoryHashSource[TH]) Set(ctx context.Context, hash TH, value []byte) error {
	a.Lock()
	a.indexes[hash] = value
	a.Unlock()
	return nil
}

func (a *memoryHashSource[TH]) Get(ctx context.Context, hash TH) ([]byte, error) {
	a.RLock()
	res, ok := a.indexes[hash]
	a.RUnlock()

	if ok {
		return res, nil
	}

	return res, types.ErrKeyNotFound
}

func (a *memoryHashSource[TH]) Delete(ctx context.Context, hash TH) error {
	a.Lock()
	delete(a.indexes, hash)
	a.Unlock()
	return nil
}
