package store

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

type IIndexSource[K index.IndexValue, V types.HashType] interface {
	GetHash(ctx context.Context, isLeaf bool, index K) (V, error)
	SetHash(ctx context.Context, isLeaf bool, index K, value V) error
}

type IHashSource[K types.HashType] interface {
	Get(ctx context.Context, hash K) ([]byte, error)
	Set(ctx context.Context, hash K, value []byte) error
	Delete(ctx context.Context, hash K) error
}
