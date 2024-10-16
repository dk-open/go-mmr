package merkle_test

import (
	"context"
	"fmt"
	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"github.com/dk-open/go-mmr/types/hasher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMmrWithDifferentHashers(t *testing.T) {
	// Create a context for the test
	ctx := context.Background()

	var numLeaves int = 10
	// Test with hasher.Ripemd160
	t.Run("Test with Ripemd160 Hashing", func(t *testing.T) {
		testMmrWithHasher[uint64, types.Hash160](t, ctx, hasher.Ripemd160, numLeaves)
	})

	// Test with hasher.Argon2 using types.Hash256
	t.Run("Test with Argon2 Hashing", func(t *testing.T) {
		testMmrWithHasher[uint64, types.Hash256](t, ctx, hasher.Argon2, numLeaves)
	})

	// Test with hasher.Blake3 using types.Hash256
	t.Run("Test with Blake3 Hashing", func(t *testing.T) {
		testMmrWithHasher[uint64, types.Hash256](t, ctx, hasher.Blake3, numLeaves)
	})
}

func TestMmrProof(t *testing.T) {
	ctx := context.Background()
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()
	m := merkle.NewMountainRange[uint64, types.Hash256](hasher.Sha3_256, memoryIndexes)

	var hashes []types.Hash256
	for i := 0; i < 7; i++ {
		data := []byte(fmt.Sprintf("test data %d", i))
		h := hasher.Sha3_256(data)
		hashes = append(hashes, h)
	}

	fmt.Println("Adding", len(hashes))
	if err := m.Add(ctx, hashes...); err != nil {
		t.Fatalf("failed to add hashes %d: %v", len(hashes), err)
	}

	p, err := m.Proof(ctx, hashes[3])
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	root, err := m.Root(ctx)
	assert.NoError(t, err, "failed to get the root of the MMR")
	assert.True(t, root.ValidateProof(p))
}

func TestMmrProofByIndex(t *testing.T) {
	ctx := context.Background()
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()
	m := merkle.NewMountainRange[uint64, types.Hash256](hasher.Sha3_256, memoryIndexes)

	for i := 0; i < 7; i++ {
		data := []byte(fmt.Sprintf("test data %d", i))

		h := hasher.Sha3_256(data)
		fmt.Println("Adding", h)
		if err := m.Add(ctx, h); err != nil {
			t.Fatalf("failed to add hash %v at index %d: %v", h, i, err)
		}
	}
	p, err := m.ProofByIndex(ctx, 3)
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	root, err := m.Root(ctx)
	assert.NoError(t, err, "failed to get the root of the MMR")
	assert.True(t, root.ValidateProof(p))
}

func TestCreateAndValidateProof_DifferentMMRSizes(t *testing.T) {
	ctx := context.Background()

	// Define multiple hashers and index types for different test cases
	testCases := []struct {
		name      string
		hasher    types.Hasher[types.Hash256]
		indexType string
		mmrSize   int
	}{
		{"Test with SHA256, uint64, MMR size 1", hasher.Sha256, "uint64", 1},
		{"Test with SHA256, uint64, MMR size 2", hasher.Sha256, "uint64", 2},
		{"Test with SHA256, uint64, MMR size 3", hasher.Sha256, "uint64", 3},
		{"Test with SHA256, uint64, MMR size 4", hasher.Sha256, "uint64", 4},
		{"Test with Blake2b, uint64, MMR size 5", hasher.Blake3, "uint64", 5},
		{"Test with Blake2b, uint64, MMR size 6", hasher.Argon2, "uint64", 6},
		{"Test with Blake2b, uint64, MMR size 7", hasher.Argon2, "uint64", 7},
		{"Test with Blake2b, uint64, MMR size 8", hasher.Argon2, "uint64", 8},
		{"Test with Blake2b, uint64, MMR size 9", hasher.Blake2b_256, "uint64", 9},
		{"Test with Blake2b, uint64, MMR size 10", hasher.Blake2b_256, "uint64", 10},
		{"Test with Blake2b, uint64, MMR size 11", hasher.Blake2b_256, "uint64", 11},
		{"Test with Blake2b, uint64, MMR size 12", hasher.Blake2b_256, "uint64", 12},
		{"Test with Blake2b, uint64, MMR size 13", hasher.Blake2b_256, "uint64", 13},
		{"Test with Blake2b, uint64, MMR size 14", hasher.Blake3, "uint64", 14},
		{"Test with Blake2b, uint64, MMR size 15", hasher.Blake3, "uint64", 15},
		{"Test with Blake2b, uint64, MMR size 16", hasher.Blake3, "uint64", 16},
		{"Test with Blake2b, uint64, MMR size 17", hasher.Blake3, "uint64", 17},
		{"Test with Blake2b, uint64, MMR size 18", hasher.Blake3, "uint64", 18},
		{"Test with Blake2b, uint64, MMR size 33", hasher.Blake3, "uint64", 33},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create an in-memory store for indexes and hashes
			memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()

			// Initialize a new MMR with the in-memory store and the selected hasher
			mmr := merkle.NewMountainRange[uint64, types.Hash256](tc.hasher, memoryIndexes)

			// Add elements to the MMR based on the specified mmrSize
			for i := 0; i < tc.mmrSize; i++ {
				data := []byte(fmt.Sprintf("Test text %d", i)) // Create sample data for each element
				hash := tc.hasher(data)
				err := mmr.Add(ctx, hash)
				assert.NoError(t, err, "should add the element without error")
			}

			// Test Proof Creation for index 0 and last index in MMR (to test edge cases)
			var i int
			for i < tc.mmrSize-1 {
				proof, err := mmr.ProofByIndex(ctx, uint64(i))
				assert.NoError(t, err, "proof creation should not return an error")
				assert.NotNil(t, proof, "proof should not be nil")
				assert.Greater(t, len(proof.Hashes), 0, "proof should contain hashes")

				// Validate the proof using the root from the MMR
				root, err := mmr.Root(ctx)
				if !assert.NoError(t, err, "failed to get the root of the MMR") {
					t.Fatal(err)
				}
				isValid := root.ValidateProof(proof)
				assert.True(t, isValid, "the proof should be valid")
				i++
			}
		})
	}
}

func testMmrWithHasher[TIndex index.Value, THash types.HashType](t *testing.T, ctx context.Context, hf func(...[]byte) THash, numLeaves int) {
	// Initialize memory-based index and hash sources
	memoryIndexes := store.MemoryIndexSource[TIndex, THash]()

	// Create a new Merkle Mountain Range using the provided hash function
	m := merkle.NewMountainRange[TIndex, THash](hf, memoryIndexes)

	leafsData := map[TIndex]THash{}

	var ni TIndex
	t.Run("Add and Retrieve Leaves", func(t *testing.T) {
		// Adding leaves to the MMR
		for i := 0; i < numLeaves; i++ {
			data := []byte(fmt.Sprintf("test data %d", ni))
			h := hf(data)
			leafsData[ni] = h

			// Add the hash to the MMR and check for errors
			if err := m.Add(ctx, h); err != nil {
				t.Fatalf("failed to add hash %v at index %d: %v", h, ni, err)
			}
			ni += 1
		}

		// Check if all leaves were added and can be retrieved
		for i, expectedHash := range leafsData {
			retrievedHash, err := m.Get(ctx, i)
			assert.NoError(t, err, "failed to retrieve hash at index %v", i)
			assert.Equal(t, expectedHash, retrievedHash, "hash mismatch at index %v", i)
		}

		// Verify the MMR size is correct
		assert.Equal(t, uint64(numLeaves), m.Size(), "MMR size mismatch")
	})

	t.Run("Retrieve Non-Existent Leaf", func(t *testing.T) {
		_, err := m.Get(ctx, ni) // Request an index out of range
		assert.Error(t, err, "expected an error when retrieving non-existent leaf")
	})
}
