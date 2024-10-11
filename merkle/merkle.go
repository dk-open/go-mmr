package merkle

import (
	"context"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"sync"
)

type IMountainRange[TIndex types.IndexValue, THash types.HashType] interface {
	Add(ctx context.Context, value THash) error
	Get(ctx context.Context, index TIndex) (THash, error)
	Size() TIndex
}

type mmr[TIndex types.IndexValue, THash types.HashType] struct {
	sync.RWMutex
	root    IRoot[TIndex, THash]
	hf      types.Hasher[THash]
	indexes store.IIndexSource[TIndex, THash]
	hashes  store.IHashSource[THash]
}

// NewMountainRange creates a new Merkle Mountain Range.
func NewMountainRange[TIndex types.IndexValue, THash types.HashType](hf types.Hasher[THash], indexes store.IIndexSource[TIndex, THash], hashes store.IHashSource[THash]) IMountainRange[TIndex, THash] {
	return &mmr[TIndex, THash]{
		indexes: indexes,
		hf:      hf,
		root:    newRoot[TIndex, THash](hf),
		hashes:  hashes,
	}
}

func (m *mmr[TIndex, THash]) Get(ctx context.Context, index TIndex) (res THash, err error) {
	m.RLock()
	res, err = m.indexes.GetHash(ctx, true, index)
	m.RUnlock()
	return
}

func (m *mmr[TIndex, THash]) Add(ctx context.Context, value THash) error {
	m.Lock()
	err := m.appendMerkle(ctx, value)
	m.Unlock()
	return err
}

func (m *mmr[TIndex, THash]) Size() TIndex {
	return m.root.Size()
}
