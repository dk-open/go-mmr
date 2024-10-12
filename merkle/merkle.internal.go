package merkle

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

func (m *mmr[TIndex, THash]) saveLeaf(ctx context.Context, index TIndex, value THash) error {
	if err := m.indexes.SetHash(ctx, true, index, value); err != nil {
		return err
	}
	ln := LeafNode[TIndex, THash](index)
	data, err := ln.MarshalBinary()
	if err != nil {
		return err
	}
	return m.hashes.Set(ctx, value, data)
}

func (m *mmr[TIndex, THash]) updateNode(ctx context.Context, i types.Index[TIndex], value THash) error {
	upper := i.RightUp()

	if upper == nil {
		return nil
	}
	upperNode := Node[TIndex, THash](upper.Index())
	if i.IsRight() {
		upperNode.SetRight(value)
	} else {
		upperNode.SetLeft(value)
	}

	siblingIndex := i.GetSibling()
	if sibHash, he := m.indexes.GetHash(ctx, i.IsLeaf(), siblingIndex.Index()); he == nil {
		if siblingIndex.IsRight() {
			upperNode.SetRight(sibHash)
		} else {
			upperNode.SetLeft(sibHash)
		}
	}

	packed, err := upperNode.MarshalBinary()
	if err != nil {
		return err
	}
	nodeHash := m.hf(packed)
	if err = m.hashes.Set(ctx, nodeHash, packed); err != nil {
		return err
	}

	if err = m.indexes.SetHash(ctx, false, upper.Index(), nodeHash); err != nil {
		return err
	}

	return m.updateNode(ctx, upper, nodeHash)
}

func (m *mmr[TIndex, THash]) appendMerkle(ctx context.Context, value THash) (err error) {
	nextIndex := m.root.Size()
	if err = m.saveLeaf(ctx, nextIndex, value); err != nil {
		return err
	}

	leafIndex := index.LeafIndex[TIndex](nextIndex)
	if err = m.updateNode(ctx, leafIndex, value); err != nil {
		return err
	}

	peaks := index.GetPeaks(leafIndex)
	hashes := make([]THash, len(peaks))
	for i, p := range peaks {
		h, hErr := m.indexes.GetHash(ctx, p.IsLeaf(), p.Index())
		if hErr != nil {
			return hErr
		}
		hashes[i] = h
	}
	return m.root.Increment(hashes...)
}
