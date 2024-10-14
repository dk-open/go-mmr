package store

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
	"sync"
)

type memoryIndexSource[K index.Value, V types.HashType] struct {
	sync.RWMutex
	leafs map[K]V
	nodes map[K]V
}

func (a *memoryIndexSource[K, V]) Set(ctx context.Context, isLeaf bool, index K, value V) error {
	a.Lock()
	if isLeaf {
		a.leafs[index] = value
	} else {
		a.nodes[index] = value
	}
	a.Unlock()
	return nil
}

func (a *memoryIndexSource[K, V]) Get(ctx context.Context, isLeaf bool, index K) (V, error) {
	var res V
	var ok bool
	a.RLock()

	if isLeaf {
		res, ok = a.leafs[index]
	} else {
		res, ok = a.nodes[index]
	}
	a.RUnlock()
	if !ok {
		return res, types.ErrKeyNotFound
	}
	return res, nil
}

func (a *memoryIndexSource[K, V]) LeafIndex(ctx context.Context, leaf V) (res K, err error) {
	a.RLock()
	defer a.RUnlock()
	for k, v := range a.leafs {
		if v == leaf {
			return k, nil
		}
	}
	return res, types.ErrKeyNotFound
}

func MemoryIndexSource[K index.Value, V types.HashType]() IIndexSource[K, V] {
	return &memoryIndexSource[K, V]{
		leafs: make(map[K]V),
		nodes: make(map[K]V),
	}
}
