package merkle

import (
	"bytes"
	"encoding"
	"github.com/dk-open/go-mmr/types"
)

type INode[THash types.HashType] interface {
	encoding.BinaryMarshaler
	Children() []THash
	SetLeft(h THash)
	SetRight(h THash)
}

// Node is a node in the Merkle Mountain Range.
type node[THash types.HashType] [2]THash

func (n *node[THash]) Children() []THash {
	return n[:]
}

func NodeFromBinary[THash types.HashType](data []byte) (INode[THash], error) {
	n := &node[THash]{}
	if err := n.UnmarshalBinary(data); err != nil {
		return nil, err
	}

	return n, nil
}

func Node[THash types.HashType](children ...THash) INode[THash] {
	res := &node[THash]{}
	for i, ch := range children {
		res[i] = ch
	}
	return res
}

func (n *node[THash]) SetLeft(h THash) {
	n[0] = h
}

func (n *node[THash]) SetRight(h THash) {
	n[1] = h
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (n *node[THash]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	for _, ch := range n {
		err := types.BufferWrite(&buf, ch)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (n *node[THash]) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	buf := bytes.NewReader(data)
	left, err := types.BufferRead[THash](buf)
	if err != nil {
		return err
	}
	n[0] = left
	if buf.Size() > 0 {
		right, rErr := types.BufferRead[THash](buf)
		if rErr != nil {
			return rErr
		}
		n[1] = right
	}
	return nil
}
