package index

import (
	"github.com/dk-open/go-mmr/types"
)

// Node is a node in the Merkle Mountain Range.
type Node[TIndex types.IndexValue, THash types.HashType] struct {
	Index TIndex
	Child []THash
}

func (n Node[TIndex, THash]) SetLeft(h THash) Node[TIndex, THash] {
	if len(n.Child) == 0 {
		return Node[TIndex, THash]{
			Index: n.Index,
			Child: []THash{h},
		}
	}
	n.Child[0] = h
	return n
}

func (n Node[TIndex, THash]) SetRight(h THash) Node[TIndex, THash] {
	if len(n.Child) == 0 {
		var leftKey THash
		return Node[TIndex, THash]{Index: n.Index, Child: []THash{leftKey, h}}
	}
	if len(n.Child) == 1 {
		n.Child = append(n.Child, h)
		return n
	}
	n.Child[1] = h
	return n
}

func (n Node[TIndex, THash]) Left() (THash, bool) {
	if len(n.Child) > 0 {
		return n.Child[0], true
	}
	var key THash
	return key, false
}

func (n Node[TIndex, THash]) Right() (THash, bool) {
	if len(n.Child) > 1 {
		return n.Child[1], true
	}
	var key THash
	return key, false
}

func HashBytes[TH types.HashType](value TH) ([]byte, error) {
	switch v := any(value).(type) {
	case types.Hash128:
		return v[:], nil
	case types.Hash160:
		return v[:], nil
	case types.Hash224:
		return v[:], nil
	case types.Hash256:
		return v[:], nil
	case types.Hash384:
		return v[:], nil
	case types.Hash512:
		return v[:], nil
	case uint8:
		return []byte{byte(v)}, nil
	case uint16:
		return []byte{byte(v >> 8), byte(v)}, nil
	case uint32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case uint64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int8:
		return []byte{byte(v)}, nil
	case int16:
		return []byte{byte(v >> 8), byte(v)}, nil
	case int32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case uint:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case string:
		return []byte(v), nil
	default:
		return nil, types.ErrTypeMismatch
	}
}
