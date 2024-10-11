package merkle

import (
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

type IRoot[TI types.IndexValue, TH types.HashType] interface {
	Hash() TH
	Size() TI
	Increment(peaks ...TH) error
}

type root[TI types.IndexValue, TH types.HashType] struct {
	index TI
	hash  TH
	hf    types.Hasher[TH]
	peaks []TH
}

func (r *root[TI, TH]) Increment(peaks ...TH) error {
	var hashes [][]byte
	for _, peak := range peaks {
		data, err := index.HashBytes[TH](peak)
		if err != nil {
			return err
		}
		hashes = append(hashes, data)
	}
	r.hash = r.hf(hashes...)
	r.index = types.AddInt(r.index, 1)
	return nil
}

func (r *root[TI, TH]) Hash() TH {
	return r.hash
}

func (r *root[TI, TH]) Size() TI {
	return r.index
}

func newRoot[TI types.IndexValue, TH types.HashType](hf types.Hasher[TH]) IRoot[TI, TH] {
	res := &root[TI, TH]{
		hf: hf,
	}
	return res
}
