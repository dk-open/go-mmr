package go_mmr

import (
	"context"
	"fmt"
	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/dk-open/go-mmr/store"
	"github.com/dk-open/go-mmr/types"
	"github.com/dk-open/go-mmr/types/hasher"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func TestMmrProof(t *testing.T) {
	ctx := context.Background()
	memoryIndexes := store.MemoryIndexSource[uint64, types.Hash256]()
	m := merkle.NewMountainRange[uint64, types.Hash256](hasher.Sha3_256, memoryIndexes)

	var hashes []types.Hash256
	for i := 0; i < 10; i++ {
		data := []byte(fmt.Sprintf("test data %d", i))

		h := hasher.Sha3_256(data)
		hashes = append(hashes, h)
	}
	fmt.Printf("Adding %d hashes\n", len(hashes))
	if err := m.Add(ctx, hashes...); err != nil {
		t.Fatalf("failed to add hash %d %v", hashes, err)
	}

	p, err := m.ProofByIndex(ctx, 3)
	if err != nil {
		t.Fatalf("failed to create proof: %v", err)
	}

	root, err := m.Root(ctx)
	assert.NoError(t, err, "failed to get the root of the MMR")
	assert.True(t, root.ValidateProof(p))
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

func benchmarkMmrWithHasher[TIndex index.Value, THash types.HashType](b *testing.B, ctx context.Context, hf func(...[]byte) THash, numElements int) {
	// Initialize memory-based index and hash sources
	memoryIndexes := store.MemoryIndexSource[TIndex, THash]()

	// Create a new Merkle Mountain Range using the provided hash function
	m := merkle.NewMountainRange[TIndex, THash](hf, memoryIndexes)

	var ni TIndex

	// Run the benchmark
	for n := 0; n < b.N; n++ {
		hashes := []THash{}
		// Add 100k elements to the MMR
		for i := 0; i < numElements; i++ {
			data := []byte(fmt.Sprintf("test data %d", ni))
			hashes = append(hashes, hf(data))

			// Increment index for each element
			ni += 1
		}
		// Add the hash to the MMR
		if err := m.Add(ctx, hashes...); err != nil {
			b.Fatalf("failed to add hash %d at index %d: %v", len(hashes), ni, err)
		}
	}
	b.ReportAllocs()
}
