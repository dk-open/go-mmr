package merkle

import (
	"encoding/binary"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
	"math/big"
)

func packIndex[TI types.IndexValue](index TI) ([]byte, error) {
	var res []byte
	switch v := any(index).(type) {
	case uint:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(res, uint64(any(index).(uint)))
		return res[:n], nil
	case uint32:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(res, uint64(any(index).(uint32)))
		return res[:n], nil
	case uint64:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(res, any(index).(uint64))
		return res[:n], nil
	case int:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(res, any(index).(int64))
		return res[:n], nil
	case int32:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(res, int64(any(index).(int32)))
		return res[:n], nil
	case int64:
		res = make([]byte, binary.MaxVarintLen64)
		n := binary.PutVarint(res, any(index).(int64))
		return res[:n], nil
	case *big.Int:
		res = make([]byte, binary.MaxVarintLen64)
		data := v.Bytes()
		n := binary.PutUvarint(res, uint64(len(data)))
		res = append(res[:n], v.Bytes()...)
	default:
		return nil, types.ErrTypeMismatch
	}
	return res, nil
}

func packNode[TI types.IndexValue, TH types.HashType](node index.Node[TI, TH]) ([]byte, error) {
	data, err := packIndex[TI](node.Index)
	if err != nil {
		return nil, err
	}
	for _, child := range node.Child {
		switch any(child).(type) {
		case types.Hash160:
			h := any(child).(types.Hash160)
			data = append(data, h[:]...)
		case types.Hash256:
			h := any(child).(types.Hash256)
			data = append(data, h[:]...)
		case types.Hash512:
			h := any(child).(types.Hash512)
			data = append(data, h[:]...)
		case types.Hash128:
			h := any(child).(types.Hash128)
			data = append(data, h[:]...)
		case types.Hash224:
			h := any(child).(types.Hash224)
			data = append(data, h[:]...)
		case types.Hash384:
			h := any(child).(types.Hash384)
			data = append(data, h[:]...)
		default:
			return nil, types.ErrTypeMismatch
		}

	}
	return data, nil
}

func unpackIndex[TI types.IndexValue](data []byte) (TI, []byte, error) {
	//index, n := binary.Uvarint(data)
	var res TI
	switch any(res).(type) {
	case uint:
		index, n := binary.Uvarint(data)
		data = data[n:]
		res = any(uint(index)).(TI)
	case uint32:
		index, n := binary.Uvarint(data)
		data = data[n:]
		res = any(uint32(index)).(TI)
	case uint64:
		index, n := binary.Uvarint(data)
		data = data[n:]
		res = any(index).(TI)
	case int:
		index, n := binary.Varint(data)
		data = data[n:]
		res = any(int(index)).(TI)
	case int32:
		index, n := binary.Varint(data)
		data = data[n:]
		res = any(int32(index)).(TI)
	case int64:
		index, n := binary.Varint(data)
		data = data[n:]
		res = any(index).(TI)
	case *big.Int:

		size, n := binary.Uvarint(data)
		data = data[n:]
		res = any(new(big.Int).SetBytes(data[:size])).(TI)
		data = data[size:]
	default:
		return res, nil, types.ErrTypeMismatch
	}

	return res, data, nil
}

func unpackNode[TI types.IndexValue, TH types.HashType](data []byte) (res index.Node[TI, TH], err error) {
	var th TH
	res.Index, data, err = unpackIndex[TI](data)
	if err != nil {
		return res, err
	}

	switch any(th).(type) {
	case types.Hash160:
		var left, right types.Hash160
		copy(left[:], data[:20])
		copy(right[:], data[20:40])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	case types.Hash256:
		var left, right types.Hash256
		copy(left[:], data[:32])
		copy(right[:], data[32:64])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	case types.Hash512:
		var left, right types.Hash512
		copy(left[:], data[:64])
		copy(right[:], data[64:128])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	case types.Hash128:
		var left, right types.Hash128
		copy(left[:], data[:16])
		copy(right[:], data[16:32])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	case types.Hash224:
		var left, right types.Hash224
		copy(left[:], data[:28])
		copy(right[:], data[28:56])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	case types.Hash384:
		var left, right types.Hash384
		copy(left[:], data[:48])
		copy(right[:], data[48:96])
		res.Child = []TH{any(left).(TH), any(right).(TH)}
	default:
		return res, types.ErrTypeMismatch
	}
	return res, nil
}
