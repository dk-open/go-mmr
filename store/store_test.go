package store_test

import (
	"context"
	"github.com/dk-open/go-mmr/store"
	"testing"

	"github.com/dk-open/go-mmr/types"
	"github.com/stretchr/testify/assert"
)

func TestMemoryIndexSource_SetAndGet(t *testing.T) {
	ctx := context.Background()

	// Create a new MemoryIndexSource instance
	source := store.MemoryIndexSource[uint32, types.Hash256]()

	// Define some mock data
	var index1 uint32 = 1
	var index2 uint32 = 2
	value1 := types.Hash256{1, 2, 3}
	value2 := types.Hash256{4, 5, 6}

	// Test setting and getting leaf values
	err := source.Set(ctx, true, index1, value1)
	assert.NoError(t, err, "Set leaf value should not return an error")

	res, err := source.Get(ctx, true, index1)
	assert.NoError(t, err, "Get leaf value should not return an error")
	assert.Equal(t, value1, res, "Get should return the correct leaf value")

	// Test setting and getting node values
	err = source.Set(ctx, false, index2, value2)
	assert.NoError(t, err, "Set node value should not return an error")

	res, err = source.Get(ctx, false, index2)
	assert.NoError(t, err, "Get node value should not return an error")
	assert.Equal(t, value2, res, "Get should return the correct node value")

	leafIndex, err := source.LeafIndex(ctx, value1)
	assert.NoError(t, err, "LeafIndex should not return an error")
	assert.Equal(t, index1, leafIndex, "LeafIndex should return the correct index")
}

func TestMemoryIndexSource_KeyNotFound(t *testing.T) {
	ctx := context.Background()

	// Create a new MemoryIndexSource instance
	source := store.MemoryIndexSource[uint64, types.Hash256]()

	// Test getting a non-existent leaf value
	_, err := source.Get(ctx, true, uint64(999))
	assert.ErrorIs(t, err, types.ErrKeyNotFound, "Getting a non-existent leaf value should return ErrKeyNotFound")

	// Test getting a non-existent node value
	_, err = source.Get(ctx, false, uint64(999))
	assert.ErrorIs(t, err, types.ErrKeyNotFound, "Getting a non-existent node value should return ErrKeyNotFound")
}

func TestMemoryIndexSource_OverwriteValue(t *testing.T) {
	ctx := context.Background()

	// Create a new MemoryIndexSource instance
	source := store.MemoryIndexSource[uint64, types.Hash256]()

	// Define some mock data
	index1 := uint64(1)
	value1 := types.Hash256{1, 2, 3}
	value2 := types.Hash256{7, 8, 9}

	// Set a leaf value
	err := source.Set(ctx, true, index1, value1)
	assert.NoError(t, err, "Set leaf value should not return an error")

	// Overwrite the leaf value
	err = source.Set(ctx, true, index1, value2)
	assert.NoError(t, err, "Set (overwrite) leaf value should not return an error")

	// Verify that the value was overwritten
	res, err := source.Get(ctx, true, index1)
	assert.NoError(t, err, "Get leaf value should not return an error")
	assert.Equal(t, value2, res, "Get should return the overwritten leaf value")
}
