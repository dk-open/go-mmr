package index

import (
	"github.com/dk-open/go-mmr/types"
	"math/big"
)

// GetPeaks Calculates Peaks
// Algorithm:
//  1. Get Top from the current position
//  2. Go To the left branch.
//     - if No Any left branches - return
//     - go To 1
func GetPeaks[TI types.IndexValue](x types.Index[TI]) (res []types.Index[TI]) {
	var peak = x
	for {
		peak = peak.Top()
		res = append(res, peak)
		if peak = peak.LeftBranch(); peak == nil {
			return
		}
	}
}

func isLeft[IV types.IndexValue](value IV) bool {
	return !isRight(value)
}

func isRight[IV types.IndexValue](value IV) bool {
	switch v := any(value).(type) {
	case int:
		return firstBitSet(v)
	case int32:
		return firstBitSet(v)
	case int64:
		return firstBitSet(v)
	case uint:
		return firstBitSet(v)
	case uint32:
		return firstBitSet(v)
	case uint64:
		return firstBitSet(v)
	case *big.Int:
		return false
	default:
		panic("unsupported type")
	}
}

func getHeight[IV types.IndexValue](value IV) (height int) {
	switch v := any(value).(type) {
	case int:
		return getNumericHeight(v)
	case int32:
		return getNumericHeight(v)
	case int64:
		return getNumericHeight(v)
	case uint:
		return getNumericHeight(v)
	case uint32:
		return getNumericHeight(v)
	case uint64:
		return getNumericHeight(v)
	case *big.Int:
		heightBig := big.NewInt(0)
		one := big.NewInt(1)
		for v.Cmp(big.NewInt(0)) != 0 && v.Bit(0) == 0 {
			v.Rsh(v, 1)
			heightBig.Add(heightBig, one)
		}
		return int(heightBig.Int64())
	default:
		panic("unsupported type")
	}
}

// getNumericHeight calculates the height of a node in a binary tree or MMR (Merkle Mountain Range).
// The height is determined by how many times the node's index can be divided by 2 (right-shifted) before it becomes odd.
// Visualization:
//
//	         [4] (height 2)
//	          |
//	 [2] (height 1)      [6] (height 1)
//	/   \                /    \
//
// [1]     [3]            [5]    [7] (height 0)
// / \      / \           / \    /  \
// 0   1    2   3       4   5   6    7
// Parameters:
// - value: The index of the node (expected to implement types.NumericValue).
//
// Returns:
// - The height of the node, as an integer, where 0 is a leaf, 1 is one level up, and so on.
func getNumericHeight[T types.NumericValue](value T) (height int) {
	var v1 = value
	for v1 != 0 && v1&1 == 0 {
		v1 = v1 >> 1
		height++
	}
	return height
}

func firstBitSet[T types.NumericValue](value T) bool {
	return value&1 == 1
}
