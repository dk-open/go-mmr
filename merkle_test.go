package go_mmr

import (
	"context"
	"fmt"
	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"github.com/dk-open/go-mmr/types/hasher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMmrProof(t *testing.T) {
	ctx := context.Background()
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()
	m := merkle.NewMountainRange[uint64, types.Hash256](hasher.Sha3_256, memoryIndexes)

	for i := 0; i < 10; i++ {
		data := []byte(fmt.Sprintf("test data %d", i))

		h := hasher.Sha3_256(data)
		fmt.Printf("Adding %x\n", h)
		if err := m.Add(ctx, h); err != nil {
			t.Fatalf("failed to add hash %v at index %d: %v", h, i, err)
		}
	}
	p, err := m.CreateProof(ctx, 3)
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	root := m.Root()
	assert.True(t, root.ValidateProof(p))
}
