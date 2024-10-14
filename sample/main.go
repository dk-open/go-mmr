//go:build !test

package main

import (
	"context"
	"fmt"
	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"github.com/dk-open/go-mmr/types/hasher"
	"log"
)

func main() {
	ctx := context.Background()
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()
	m := merkle.NewMountainRange[uint64, types.Hash256](hasher.Sha3_256, memoryIndexes)

	for i := 0; i < 10; i++ {
		h := hasher.Sha3_256([]byte(fmt.Sprintf("test data %d", i)))
		fmt.Printf("Adding at %d item %x\n", i, h)
		if err := m.Add(ctx, h); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println()
	root := m.Root()
	fmt.Printf("Root: %x\n", root)
	item3, err := m.Get(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Item at index 3: %x\n", item3)
	fmt.Println("Create a proof for item 4")
	prooft, err := m.CreateProof(ctx, 4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Validate a proof")
	if !root.ValidateProof(prooft) {
		log.Fatal("Proof is not invalid")
	}

	fmt.Println("Proof is valid")
}
