package store

import (
	"context"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

type IIndexSource[K index.Value, V types.HashType] interface {
	Get(ctx context.Context, isLeaf bool, index K) (V, error)
	Set(ctx context.Context, isLeaf bool, index K, value V) error
	LeafIndex(ctx context.Context, leaf V) (K, error)
}
