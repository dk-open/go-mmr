package merkle

import (
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/types"
)

type Proof[TIndex index.IndexValue, THash types.HashType] struct {
	Target     TIndex
	Hashes     []THash
	LeftPeaks  []THash
	RightPeaks []THash
}
