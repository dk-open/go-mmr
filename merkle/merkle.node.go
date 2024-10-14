package merkle

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
	"io"
	"math/big"
)

type INode[TIndex index.IndexValue, THash types.HashType] interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	IsLeaf() bool
	Index() TIndex
	Children() []THash
	Left() (res THash, ok bool)
	Right() (res THash, ok bool)
	SetLeft(h THash) INode[TIndex, THash]
	SetRight(h THash) INode[TIndex, THash]
}

// Node is a node in the Merkle Mountain Range.
type node[TIndex index.IndexValue, THash types.HashType] struct {
	leaf     bool
	index    TIndex
	children []THash
}

func (n *node[TIndex, THash]) Index() TIndex {
	return n.index
}

func (n *node[TIndex, THash]) IsLeaf() bool {
	return n.leaf
}

func (n *node[TIndex, THash]) Children() []THash {
	return n.children
}

func LeafNode[TIndex index.IndexValue, THash types.HashType](index TIndex) INode[TIndex, THash] {
	return &node[TIndex, THash]{
		leaf:  true,
		index: index,
	}
}

func NodeFromBinary[TIndex index.IndexValue, THash types.HashType](data []byte) (INode[TIndex, THash], error) {
	n := &node[TIndex, THash]{}
	if err := n.UnmarshalBinary(data); err != nil {
		return nil, err
	}

	return n, nil
}

func Node[TIndex index.IndexValue, THash types.HashType](index TIndex, children ...THash) INode[TIndex, THash] {
	return &node[TIndex, THash]{
		leaf:     true,
		index:    index,
		children: children,
	}
}

func (n *node[TIndex, THash]) SetLeft(h THash) INode[TIndex, THash] {
	if len(n.children) == 0 {
		return &node[TIndex, THash]{
			index:    n.index,
			children: []THash{h},
		}
	}
	n.children[0] = h
	return n
}

func (n *node[TIndex, THash]) SetRight(h THash) INode[TIndex, THash] {
	if len(n.children) == 0 {
		var leftKey THash
		return &node[TIndex, THash]{
			index:    n.index,
			leaf:     n.leaf,
			children: []THash{leftKey, h},
		}
	}
	if len(n.children) == 1 {
		n.children = append(n.children, h)
		return n
	}
	n.children[1] = h
	return n
}

func (n *node[TIndex, THash]) Left() (res THash, ok bool) {
	if len(n.children) > 0 {
		return n.children[0], true
	}
	return
}

func (n *node[TIndex, THash]) Right() (res THash, ok bool) {
	if len(n.children) > 1 {
		return n.children[1], true
	}
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (n *node[TIndex, THash]) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	// Marshal the 'leaf' field (bool to byte)
	if n.leaf {
		buf.WriteByte(1)
	} else {
		buf.WriteByte(0)
	}

	// Marshal the 'index' field using variant type encoding
	if err := encodeVariantIndex(&buf, n.index); err != nil {
		return nil, err
	}

	if err := encodeVariantIndex[int](&buf, len(n.children)); err != nil {
		return nil, err
	}

	for _, ch := range n.children {
		err := writeHash(&buf, ch)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (n *node[TIndex, THash]) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)

	// Unmarshal the 'leaf' field
	leafByte, err := buf.ReadByte()
	if err != nil {
		return err
	}
	n.leaf = leafByte != 0

	// Unmarshal the 'index' field using variant type decoding
	if n.index, err = decodeVariant[TIndex](buf); err != nil {
		return err
	}

	// Unmarshal the 'child' field (slice of THash)
	childLength, err := decodeVariant[int](buf)
	if err != nil {
		return err
	}
	if childLength == 0 {
		return nil
	}
	n.children = make([]THash, childLength)
	for i := 0; i < childLength; i++ {
		if n.children[i], err = readHash[THash](buf); err != nil {
			return err
		}
	}
	return nil
}

// Helper: Variant encoding for signed integers using Zigzag encoding
func zigzagEncode(n int64) uint64 {
	return uint64((n << 1) ^ (n >> 63))
}

func zigzagDecode(n uint64) int64 {
	return int64((n >> 1) ^ uint64(int64(n<<63)>>63))
}

// Helper: Encode numeric index values as a variant type
func encodeVariantIndex[TIndex index.IndexValue](buf *bytes.Buffer, index TIndex) error {
	switch v := any(index).(type) {
	case int:
		return writeUVarint(buf, uint64(v))
	case uint:
		return writeUVarint(buf, uint64(v))
	case uint32:
		return writeUVarint(buf, uint64(v))
	case uint64:
		return writeUVarint(buf, v)
	case int32:
		return writeUVarint(buf, zigzagEncode(int64(v)))
	case int64:
		return writeUVarint(buf, zigzagEncode(v))
	case *big.Int:
		b := v.Bytes()
		binary.Write(buf, binary.LittleEndian, uint64(len(b))) // Write length
		_, err := buf.Write(b)                                 // Write big.Int bytes
		return err
	default:
		return fmt.Errorf("unsupported index type for variant encoding")
	}
}

// Helper: Decode numeric index values as a variant type
func decodeVariant[TIndex index.IndexValue](r io.Reader) (TIndex, error) {
	switch any(*new(TIndex)).(type) {
	case int:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err
		}
		return any(int(v)).(TIndex), nil
	case uint:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err

		}
		return any(uint(v)).(TIndex), nil

	case uint32:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err
		}
		return any(uint32(v)).(TIndex), nil
	case uint64:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err
		}
		return any(v).(TIndex), nil
	case int32:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err
		}
		return any(int32(zigzagDecode(v))).(TIndex), nil
	case int64:
		v, err := readUVarint(r)
		if err != nil {
			return *new(TIndex), err
		}
		return any(zigzagDecode(v)).(TIndex), nil
	case *big.Int:
		var length uint64
		err := binary.Read(r, binary.LittleEndian, &length) // Read length
		if err != nil {
			var res TIndex
			return res, err
		}
		b := make([]byte, length)
		_, err = r.Read(b) // Read big.Int bytes
		if err != nil {
			var res TIndex
			return res, err
		}
		return any(new(big.Int).SetBytes(b)).(TIndex), nil
	default:
		return *new(TIndex), fmt.Errorf("unsupported index type for variant decoding")
	}
}

// Write unsigned varint (LEB128 style)
func writeUVarint(buf *bytes.Buffer, value uint64) error {
	for value >= 0x80 {
		buf.WriteByte(byte(value) | 0x80)
		value >>= 7
	}
	buf.WriteByte(byte(value))
	return nil
}

// Read unsigned varint (LEB128 style)
func readUVarint(r io.Reader) (uint64, error) {
	var result uint64
	var shift uint
	for i := 0; ; i++ {
		var b byte
		err := binary.Read(r, binary.LittleEndian, &b)
		if err != nil {
			return 0, err
		}
		result |= uint64(b&0x7F) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

// Helper to write the hash field to the buffer (same as before)
func writeHash[THash types.HashType](buf *bytes.Buffer, hash THash) error {
	switch any(hash).(type) {
	case types.Hash160:
		h := any(hash).(types.Hash160)
		buf.Write(h[:])
	case types.Hash256:
		h := any(hash).(types.Hash256)
		buf.Write(h[:])
	case types.Hash512:
		h := any(hash).(types.Hash512)
		buf.Write(h[:])
	case types.Hash128:
		h := any(hash).(types.Hash128)
		buf.Write(h[:])
	case types.Hash224:
		h := any(hash).(types.Hash224)
		buf.Write(h[:])
	case types.Hash384:
		h := any(hash).(types.Hash384)
		buf.Write(h[:])
	default:
		return types.ErrTypeMismatch
	}
	return nil
}

// Helper to read the hash field from the buffer (same as before)
func readHash[THash types.HashType](buf *bytes.Reader) (res THash, err error) {
	switch v := any(res).(type) {
	case types.Hash128:
		//var hash types.Hash128
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case types.Hash160:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case types.Hash224:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case types.Hash256:
		_, err = buf.Read(v[:32])
		return any(v).(THash), err
	case types.Hash384:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case types.Hash512:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	default:
		return res, types.ErrTypeMismatch
	}
}
