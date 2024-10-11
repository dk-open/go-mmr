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

func TestMmr(t *testing.T) {
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash160]()
	memoryHashes := store.MemoryHashSource[types.Hash160]()
	hf := hasher.Ripemd160

	m := merkle.NewMountainRange[uint64, types.Hash160](hf, memoryIndexes, memoryHashes)
	fmt.Println(m)
	ctx := context.Background()

	leafsData := map[uint64]types.Hash160{}
	i := 0
	for i < 10 {
		data := []byte(fmt.Sprintf("test data %d", i))
		h := hf(data)

		leafsData[uint64(i)] = h
		if err := m.Add(ctx, h); err != nil {
			t.Fatal(err)
			return
		}
		i++
	}

	for index, h1 := range leafsData {
		h2, err := m.Get(ctx, index)
		if err != nil {
			t.Fatal(err)
			return
		}
		assert.Equal(t, h1, h2)
	}

	assert.Equal(t, uint64(10), m.Size())
}
