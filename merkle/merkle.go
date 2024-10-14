package merkle

import (
	"context"
	"errors"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"sync"
)

type IMountainRange[TIndex index.Value, THash types.HashType] interface {
	Add(ctx context.Context, value THash) error
	Get(ctx context.Context, index TIndex) (THash, error)
	CreateProof(ctx context.Context, index TIndex) (*Proof[TIndex, THash], error)
	Root() IRoot[TIndex, THash]
	Size() TIndex
}

type mmr[TIndex index.Value, THash types.HashType] struct {
	sync.RWMutex
	root    THash
	size    TIndex
	hf      types.Hasher[THash]
	indexes store.IIndexSource[TIndex, THash]
}

// NewMountainRange creates a new Merkle Mountain Range.
func NewMountainRange[TIndex index.Value, THash types.HashType](hf types.Hasher[THash], indexes store.IIndexSource[TIndex, THash]) IMountainRange[TIndex, THash] {
	return &mmr[TIndex, THash]{
		indexes: indexes,
		hf:      hf,
	}
}

func (m *mmr[TIndex, THash]) Get(ctx context.Context, index TIndex) (res THash, err error) {
	m.RLock()
	res, err = m.indexes.Get(ctx, true, index)
	m.RUnlock()
	return
}

func (m *mmr[TIndex, THash]) Add(ctx context.Context, value THash) error {
	m.Lock()
	err := m.appendMerkle(ctx, value)
	m.Unlock()
	return err
}

// getProofIndexes collects the indexes needed to create a proof for the given item.
func (m *mmr[TIndex, THash]) getProofIndexes(item index.Index[TIndex], maxIndex TIndex) []index.Index[TIndex] {
	// Initialize the result with the item itself.
	res := make([]index.Index[TIndex], 0, 10)
	res = append(res, item)
	topIndex := item
	sibIndex := item.GetSibling()
	if sibIndex != nil && sibIndex.Index() <= maxIndex {
		topIndex = sibIndex
		for sibIndex != nil && sibIndex.Index() <= maxIndex {
			res = append(res, sibIndex)
			topIndex = sibIndex.Up()
			sibIndex = topIndex.GetSibling()
		}
	}
	return res
}

func (m *mmr[TIndex, THash]) CreateProof(ctx context.Context, i TIndex) (*Proof[TIndex, THash], error) {
	m.RLock()
	defer m.RUnlock()

	var err error
	proof := &Proof[TIndex, THash]{
		Target: i,
		Hashes: []THash{},
	}

	if i >= m.size {
		return nil, errors.New("index out of range")
	}

	peaks := index.GetPeaks[TIndex](index.LeafIndex(m.size - 1))
	var start TIndex = 0
	end := m.size
	targetPeakFound := false
	var leftPeaks []index.Index[TIndex]
	var rightPeaks []index.Index[TIndex]
	var proofIndexes []index.Index[TIndex]

	for _, p := range peaks {
		if p.IsLeaf() {
			start = p.Index()
		} else {
			start = end - (1 << (p.GetHeight() + 1))
		}

		if start <= i && i < end {
			li := index.LeafIndex(i)
			if m.size == 1 {
				proofIndexes = []index.Index[TIndex]{li}
			} else {
				proofIndexes = m.getProofIndexes(li, end)
			}
			targetPeakFound = true
		} else {
			if targetPeakFound {
				leftPeaks = append(leftPeaks, p)
			} else {
				rightPeaks = append(rightPeaks, p)
			}
		}
		end = start
	}
	proof.LeftPeaks, err = m.indexToHash(ctx, leftPeaks)
	if err != nil {
		return nil, err
	}

	proof.RightPeaks, err = m.indexToHash(ctx, rightPeaks)
	if err != nil {
		return nil, err
	}
	proof.Hashes, err = m.indexToHash(ctx, proofIndexes)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (m *mmr[TIndex, THash]) Size() TIndex {
	return m.size
}

func (m *mmr[TIndex, THash]) Root() IRoot[TIndex, THash] {
	return newRoot[TIndex, THash](m.root, m.hf)
}
