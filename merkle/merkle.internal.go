package merkle

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

func (m *mmr[TIndex, THash]) saveLeaf(ctx context.Context, i TIndex, value THash) error {
	return m.indexes.Set(ctx, true, i, value)
}

func (m *mmr[TIndex, THash]) updateNode(ctx context.Context, i index.Index[TIndex], value THash) error {
	upper := i.RightUp()

	if upper == nil {
		return nil
	}
	upperNode := Node[THash]()
	if i.IsRight() {
		upperNode.SetRight(value)
	} else {
		upperNode.SetLeft(value)
	}

	siblingIndex := i.GetSibling()
	if sibHash, he := m.indexes.Get(ctx, i.IsLeaf(), siblingIndex.Index()); he == nil {
		if siblingIndex.IsRight() {
			upperNode.SetRight(sibHash)
		} else {
			upperNode.SetLeft(sibHash)
		}
	}

	return buildNodeHash(m.hf, upperNode, func(nodeHash THash) error {
		if err := m.indexes.Set(ctx, false, upper.Index(), nodeHash); err != nil {
			return err
		}

		return m.updateNode(ctx, upper, nodeHash)
	})

}

func buildNodeHash[THash types.HashType](hf types.Hasher[THash], upperNode INode[THash], f func(THash) error) error {
	packed, err := upperNode.MarshalBinary()
	if err != nil {
		return err
	}
	nodeHash := hf(packed)
	return f(nodeHash)
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
		h, hErr := m.indexes.Get(ctx, p.IsLeaf(), p.Index())
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
		if h, err := m.indexes.Get(ctx, nodeIndex.IsLeaf(), nodeIndex.Index()); err != nil {
			return nil, err
		} else {
			res[i] = h
		}
		i++
	}
	return res, nil
}
