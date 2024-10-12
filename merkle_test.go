package go_mmr

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"

	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"github.com/dk-open/go-mmr/types/hasher"
	"github.com/stretchr/testify/assert"
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

func testMmrWithHasher[TIndex types.IndexValue, THash types.HashType](t *testing.T, ctx context.Context, hf func(...[]byte) THash, numLeaves int) {
	// Initialize memory-based index and hash sources
	memoryIndexes := store.MemoryIndexSource[TIndex, THash]()
	memoryHashes := store.MemoryHashSource[THash]()

	// Create a new Merkle Mountain Range using the provided hash function
	m := merkle.NewMountainRange[TIndex, THash](hf, memoryIndexes, memoryHashes)

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
			ni = types.AddInt(ni, 1)
		}

		// Check if all leaves were added and can be retrieved
		for index, expectedHash := range leafsData {
			retrievedHash, err := m.Get(ctx, index)
			assert.NoError(t, err, "failed to retrieve hash at index %v", index)
			assert.Equal(t, expectedHash, retrievedHash, "hash mismatch at index %v", index)
		}

		// Verify the MMR size is correct
		assert.Equal(t, uint64(numLeaves), m.Size(), "MMR size mismatch")
	})

	t.Run("Retrieve Non-Existent Leaf", func(t *testing.T) {
		_, err := m.Get(ctx, ni) // Request an index out of range
		assert.Error(t, err, "expected an error when retrieving non-existent leaf")
	})
}

func BenchmarkMmrWithDifferentHashers(b *testing.B) {
	// Create a context for the benchmark
	ctx := context.Background()
	numElements := 100 // Benchmark for 100k elements

	// Benchmark with hasher.Ripemd160
	b.Run("Benchmark with Ripemd160", func(b *testing.B) {
		benchmarkMmrWithHasher[uint64, types.Hash160](b, ctx, hasher.Ripemd160, numElements)
		b.ReportAllocs()
	})

	// Benchmark with hasher.Blake3 using types.Hash256
	b.Run("Benchmark with Blake3", func(b *testing.B) {
		benchmarkMmrWithHasher[uint64, types.Hash256](b, ctx, hasher.Blake3, numElements)
		b.ReportAllocs()
	})

	b.Run("Benchmark with Blake2b 256", func(b *testing.B) {
		benchmarkMmrWithHasher[uint64, types.Hash256](b, ctx, hasher.Blake2b_256, numElements)
		b.ReportAllocs()
	})

	b.Run("Benchmark with Blake2b 512", func(b *testing.B) {
		benchmarkMmrWithHasher[uint32, types.Hash512](b, ctx, hasher.Blake2b_512, numElements)
		b.ReportAllocs()
	})
}

func BenchmarkMmrWithDifferentHashersWithProfiler(b *testing.B) {
	// Create a context for the benchmark
	ctx := context.Background()
	numElements := 10000 // Benchmark for 10k elements

	// CPU profiling
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		b.Fatal("could not create CPU profile: ", err)
	}
	defer cpuProfile.Close()

	if err = pprof.StartCPUProfile(cpuProfile); err != nil {
		b.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	// Benchmark with hasher.Ripemd160
	b.Run("Benchmark with Ripemd160", func(b *testing.B) {
		benchmarkMmrWithHasher[uint64, types.Hash160](b, ctx, hasher.Ripemd160, numElements)
	})
	time.Sleep(2 * time.Second)

	// Memory profiling
	memProfile, err := os.Create("mem.prof")
	if err != nil {
		b.Fatal("could not create memory profile: ", err)
	}
	defer memProfile.Close()
	runtime.GC()
	runtime.MemProfileRate = 1
	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		b.Fatal("could not write memory profile: ", err)
	}
}

func benchmarkMmrWithHasher[TIndex types.IndexValue, THash types.HashType](b *testing.B, ctx context.Context, hf func(...[]byte) THash, numElements int) {
	// Initialize memory-based index and hash sources
	memoryIndexes := store.MemoryIndexSource[TIndex, THash]()
	memoryHashes := store.MemoryHashSource[THash]()

	// Create a new Merkle Mountain Range using the provided hash function
	m := merkle.NewMountainRange[TIndex, THash](hf, memoryIndexes, memoryHashes)

	var ni TIndex

	// Run the benchmark
	for n := 0; n < b.N; n++ {
		// Add 100k elements to the MMR
		for i := 0; i < numElements; i++ {
			data := []byte(fmt.Sprintf("test data %d", ni))
			h := hf(data)

			// Add the hash to the MMR
			if err := m.Add(ctx, h); err != nil {
				b.Fatalf("failed to add hash %v at index %d: %v", h, ni, err)
			}

			// Increment index for each element
			ni = types.AddInt(ni, 1)
		}
	}
	b.ReportAllocs()
}
