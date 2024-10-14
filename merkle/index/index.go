package index

import (
	"math/bits"
)

// IndexValue - Separate type for index values
type IndexValue interface {
	int | int16 | int32 | int64 | uint | uint16 | uint32 | uint64
}

// Index index navigator
type Index[TI IndexValue] interface {
	GetHeight() int
	LeftBranch() Index[TI]
	GetSibling() Index[TI]
	RightUp() Index[TI]
	Up() Index[TI]
	IsRight() bool
	Top() Index[TI]
	Index() TI
	Children() []Index[TI]
	IsLeaf() bool
	Key() string
}

// GetPeaks Calculates Peaks
// Algorithm:
//  1. Get Top from the current position
//  2. Go To the left branch.
//     - if No Any left branches - return
//     - go To 1
func GetPeaks[TI IndexValue](x Index[TI]) (res []Index[TI]) {
	res = make([]Index[TI], 0, 10)
	var peak = x
	for {
		peak = peak.Top()
		res = append(res, peak)
		if peak = peak.LeftBranch(); peak == nil {
			return
		}
	}
}

// getHeight calculates the height of a node in a binary tree or MMR (Merkle Mountain Range).
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
func getHeight[IV IndexValue](value IV) (height int) {
	if value == 0 {
		return 0
	}
	return bits.TrailingZeros(uint(value))
}
