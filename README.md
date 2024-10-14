# Merkle Mountain Range (MMR) Implementation

![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dk-open/go-mmr)](https://goreportcard.com/report/github.com/dk-open/go-mmr)
[![codecov](https://codecov.io/gh/dk-open/go-mmr/graph/badge.svg?token=0UU4GMK24V)](https://codecov.io/gh/dk-open/go-mmr)
[![Go Version](https://img.shields.io/github/go-mod/go-version/dk-open/go-mmr)](https://github.com/dk-open/go-mmr)
[![GitHub release](https://img.shields.io/github/release/dk-open/go-mmr.svg)](https://github.com/dk-open/go-mmr/releases)
[![GitHub issues](https://img.shields.io/github/issues/dk-open/go-mmr)](https://github.com/dk-open/go-mmr/issues)

## Introduction

The **Merkle Mountain Range (MMR)** is a cryptographic data structure designed for efficient, append-only data handling. It builds upon the concept of Merkle Trees by allowing the accumulation of data in multiple binary trees (mountains), ensuring that no data needs to be modified or recalculated as new data is added. This is especially useful for blockchain, timestamping, and verifiable log systems where append-only structures and proofs of inclusion are crucial.

The concept was initially proposed by Peter Todd and has been implemented in various blockchain-related projects like **OpenTimestamps** and **Grin**. Our implementation, however, improves upon the original by optimizing node indexing and navigation, making data operations faster and more efficient.

### Why MMR?
- **Efficient Appending**: MMR allows for new data to be appended without requiring all previously stored data to be available.
- **Compact Proofs**: MMR enables efficient, compact proofs of data inclusion, which makes it perfect for decentralized systems where bandwidth and storage efficiency are important.
- **Scalability**: With its append-only nature, MMR structures scale very well with growing datasets.

### Key Differences in This Implementation
Our MMR implementation introduces an enhanced indexing mechanism for faster node navigation, offering a more intuitive and efficient traversal experience. The key optimizations include:
- **Optimized Node Indexing**: Improved node indexing makes navigation across the structure significantly faster compared to the original. This reduces complexity when traversing nodes, especially in large datasets.
- **Support for Multiple Hash Functions**: Unlike some implementations, ours supports multiple cryptographic hash functions like SHA-256, Blake2b, Argon2, and more. This adds flexibility for developers to choose the hashing algorithm that best fits their needs.
- **Enhanced Proof of Inclusion**: A focus on simplifying and optimizing the generation and verification of proofs for append-only data, which is especially important for decentralized and distributed systems.

For more details on the original MMR implementation:
- [OpenTimestamps MMR Documentation](https://github.com/opentimestamps/opentimestamps-server/blob/master/doc/merkle-mountain-range.md)
- [Grin MMR Documentation](https://github.com/mimblewimble/grin/blob/master/doc/mmr.md)

### Visualization of the MMR Structure
![MMR Structure](./doc/mmr-1.png)

- **Blue nodes**: Represent the actual data objects in the structure.
- **Green nodes**: Internal nodes that support the structure by linking the data in a verifiable way.

---

## Features
- **Efficient Node Storage**: For `N` data objects, approximately `N` supporting nodes are required, making the total storage size approximately `2 * N` nodes.
- **Support for Multiple Hash Functions**: Including Argon2, SHA-256, Blake2b, and more.
- **Optimized Indexing**: Simplified navigation between nodes, making traversal and data retrieval faster than traditional implementations.
- **Append-Only Nature**: You can add new elements without having access to previously appended data, improving scalability.
- **Proof Creation & Validation**: Easily create and validate cryptographic proofs of data inclusion.

### Use Cases
- **Blockchain**: Ideal for maintaining verifiable transaction histories.
- **Timestamping**: Ensures that data existed at a specific point in time, without requiring access to all historical data.
- **Verifiable Logs**: Perfect for applications where data integrity and proof of existence are crucial.

### Proof System
Proofs in MMR provide a way to prove the inclusion of specific data in the MMR without having access to the entire dataset.

1. **Create a Proof**: For any given element, the CreateProof method generates a compact proof that includes the necessary cryptographic hashes to trace the element up the MMR tree.
2. **Validate a Proof**: The ValidateProof method allows verification of the proof by recomputing the Merkle root from the proof and comparing it to the actual root.

These methods are useful in decentralized systems where bandwidth and storage efficiency are critical, such as blockchain light clients.


---

## Requirements

- **Go Version**: 1.18+

## Installation

To include this MMR implementation in your Go project, simply run:

```bash
go get -u github.com/discretemind/mmr
```

## Example Usage

```go
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
	var transactions []types.Hash256
	for i := 0; i < 10; i++ {
		h := hasher.Sha3_256([]byte(fmt.Sprintf("test data %d", i)))
		transactions = append(transactions, h)
		fmt.Printf("Adding at %d item %x\n", i, h)

	}
	if err := m.Add(ctx, transactions...); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	root, err := m.Root(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Root: %x\n", root)
	item3, err := m.Get(ctx, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Item at index 3: %x\n", item3)
	fmt.Println("Create a proof for item 4")
	prooft, err := m.ProofByIndex(ctx, 4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Validate a proof")
	if !root.ValidateProof(prooft) {
		log.Fatal("Proof is not invalid")
	}

	fmt.Println("Proof is Valid")
}

```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

