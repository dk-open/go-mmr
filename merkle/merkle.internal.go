package merkle

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

func (m *mmr[TIndex, THash]) getNode(ctx context.Context, i TIndex) (index.Node[TIndex, THash], error) {
	h, err := m.indexes.GetHash(ctx, false, i)
	if err != nil {
		return index.Node[TIndex, THash]{}, err
	}

	data, err := m.hashes.Get(ctx, h)
	if err != nil {
		return index.Node[TIndex, THash]{}, err
	}
	return unpackNode[TIndex, THash](data)
}

func (m *mmr[TIndex, THash]) saveNode(ctx context.Context, index TIndex, node index.Node[TIndex, THash], f func(ctx context.Context, h THash) error) error {
	packedNode, err := packNode[TIndex, THash](node)
	if err != nil {
		return err
	}

	nodeHash := m.hf(packedNode)
	if err = m.hashes.Set(ctx, nodeHash, packedNode); err != nil {
		return err
	}

	return f(ctx, nodeHash)
}

func (m *mmr[TIndex, THash]) saveLeaf(ctx context.Context, index TIndex, value THash) error {
	if err := m.indexes.SetHash(ctx, true, index, value); err != nil {
		return err
	}
	return nil
}

func (m *mmr[TIndex, THash]) updateNode(ctx context.Context, i types.Index[TIndex], value THash) error {
	upper := i.RightUp()

	if upper == nil {
		return nil
	}
	upperNode := index.Node[TIndex, THash]{
		Index: upper.Index(),
	}
	upperNode.Child = []THash{value}
	siblingIndex := i.GetSibling()
	if sibHash, he := m.indexes.GetHash(ctx, i.IsLeaf(), siblingIndex.Index()); he == nil {
		upperNode.Child = append(upperNode.Child, sibHash)
	}
	packedNode, err := packNode[TIndex, THash](upperNode)
	if err != nil {
		return err
	}
	nodeHash := m.hf(packedNode)
	if err = m.hashes.Set(ctx, nodeHash, packedNode); err != nil {
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
	for _, p := range peaks {
		h, hErr := m.indexes.GetHash(ctx, p.IsLeaf(), p.Index())
		if hErr != nil {
			return hErr
		}
		hashes = append(hashes, h)
	}
	return m.root.Increment(hashes...)
}
