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

func (m *mmr[TIndex, THash]) updateNode(ctx context.Context, i index.Index[TIndex], value THash) error {
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

	return buildNodeHash(m.hf, upperNode, func(nodeHash THash, packed []byte) error {
		if err := m.hashes.Set(ctx, nodeHash, packed); err != nil {
			return err
		}

		if err := m.indexes.SetHash(ctx, false, upper.Index(), nodeHash); err != nil {
			return err
		}

		return m.updateNode(ctx, upper, nodeHash)
	})

}

func buildNodeHash[TIndex index.IndexValue, THash types.HashType](hf types.Hasher[THash], upperNode INode[TIndex, THash], f func(THash, []byte) error) error {
	packed, err := upperNode.MarshalBinary()
	if err != nil {
		return err
	}
	nodeHash := hf(packed)
	return f(nodeHash, packed)
}

func (m *mmr[TIndex, THash]) appendMerkle(ctx context.Context, value THash) (err error) {
	nextIndex := m.size
	if err = m.saveLeaf(ctx, nextIndex, value); err != nil {
		return err
	}

	leafIndex := index.LeafIndex[TIndex](nextIndex)
	if err = m.updateNode(ctx, leafIndex, value); err != nil {
		return err
	}

	peaks := index.GetPeaks(leafIndex)
	hashes := make([][]byte, len(peaks))
	for i, p := range peaks {
		h, hErr := m.indexes.GetHash(ctx, p.IsLeaf(), p.Index())
		if hErr != nil {
			return hErr
		}
		data, hErr := types.HashBytes[THash](h)
		if hErr != nil {
			return hErr
		}
		hashes[i] = data
	}
	m.root = m.hf(hashes...)
	m.size = m.size + 1
	return nil
}

func (m *mmr[TIndex, THash]) indexToHash(ctx context.Context, indexes []index.Index[TIndex]) ([]THash, error) {
	res := make([]THash, len(indexes))
	i := 0
	for _, nodeIndex := range indexes {
		if h, err := m.indexes.GetHash(ctx, nodeIndex.IsLeaf(), nodeIndex.Index()); err != nil {
			return nil, err
		} else {
			res[i] = h
		}
		i++
	}
	return res, nil
}

func (m *mmr[TIndex, THash]) GetHashNode(ctx context.Context, hash THash, f func(context.Context, INode[TIndex, THash]) error) error {
	data, err := m.hashes.Get(ctx, hash)
	if err != nil {
		return err
	}
	n, err := NodeFromBinary[TIndex, THash](data)
	if err != nil {
		return err
	}
	return f(ctx, n)
}
