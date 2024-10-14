package merkle

import (
	"github.com/dk-open/go-mmr/types/hasher"
	"testing"

	"github.com/dk-open/go-mmr/types"
	"github.com/stretchr/testify/assert"
)

func TestRoot_Hash(t *testing.T) {
	// Mock hash value
	hash := types.Hash256{1, 2, 3}

	// Create a new root instance
	r := newRoot[int64, types.Hash256](hash, hasher.Argon2)

	// Test that the Hash method returns the correct value
	assert.Equal(t, hash, r.Hash(), "Hash should return the correct hash")
}

func TestRoot_ValidateProof_Success(t *testing.T) {
	// Mock proof data
	proof := &Proof[int64, types.Hash256]{
		Target:     int64(1),
		Hashes:     []types.Hash256{{1, 2, 3}, {4, 5, 6}},
		LeftPeaks:  []types.Hash256{{7, 8, 9}},
		RightPeaks: []types.Hash256{{10, 11, 12}},
	}

	// Mock root hash and hasher
	rootHash := types.Hash256{
		120, 194, 37, 98, 96, 58, 144, 45, 182, 40, 162, 40, 37, 116, 221,
		84, 120, 21, 135, 104, 30, 4, 235, 114, 148, 255, 55, 224, 154, 6, 133, 69,
	}
	hf := hasher.Sha3_256

	// Create a new root instance
	r := newRoot[int64, types.Hash256](rootHash, hf)

	// Test successful proof validation
	assert.True(t, r.ValidateProof(proof), "ValidateProof should return true for valid proof")
}

func TestRoot_ValidateProof_Failure(t *testing.T) {
	// Mock invalid proof data
	proof := &Proof[int64, types.Hash256]{
		Target:     int64(1),
		Hashes:     []types.Hash256{{1, 2, 3}, {9, 9, 9}}, // Incorrect hash to cause failure
		LeftPeaks:  []types.Hash256{{7, 8, 9}},
		RightPeaks: []types.Hash256{{10, 11, 12}},
	}

	// Mock root hash and hasher
	rootHash := types.Hash256{1, 2, 3}
	hf := hasher.Blake2b_256

	// Create a new root instance
	r := newRoot[int64, types.Hash256](rootHash, hf)

	// Test failed proof validation
	assert.False(t, r.ValidateProof(proof), "ValidateProof should return false for invalid proof")
}

func TestRoot_ValidateProof_EmptyProof(t *testing.T) {
	// Mock empty proof data
	proof := &Proof[int64, types.Hash256]{
		Target:     int64(1),
		Hashes:     []types.Hash256{}, // Empty Hashes to simulate invalid proof
		LeftPeaks:  []types.Hash256{},
		RightPeaks: []types.Hash256{},
	}

	// Mock root hash and hasher
	rootHash := types.Hash256{1, 2, 3}

	// Create a new root instance
	r := newRoot[int64, types.Hash256](rootHash, hasher.Blake3)

	// Test failed proof validation with empty proof
	assert.False(t, r.ValidateProof(proof), "ValidateProof should return false for empty proof")
}
