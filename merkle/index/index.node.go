package index

import (
	"github.com/dk-open/go-mmr/types"
	"math"
)

type nodeIndex[TI types.IndexValue] struct {
	value TI
}

// NodeIndex creates a new node index with the given value.
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
func NodeIndex[TI types.IndexValue](value TI) types.Index[TI] {
	res := &nodeIndex[TI]{
		value: value,
	}
	return res
}

// GetHeight calculates the height of a node in a binary tree or MMR (Merkle Mountain Range).
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
func (n *nodeIndex[TI]) GetHeight() int {
	return getHeight(n.value)
}

// LeftBranch calculates the index of the left child (left branch) of the current node.
// It checks if the left child exists based on the node's height and position,
// and returns the left child's index or nil if there is no left branch.
func (n *nodeIndex[TI]) LeftBranch() types.Index[TI] {
	//pow := uint64(math.Pow(2, float64(n.GetHeight()+1)))
	pow := uint64(1) << (n.GetHeight() + 1)
	if types.IndexUint64(n.value) > pow {
		return NodeIndex[TI](types.SubtractUint64(n.value, pow))
	}
	return nil
}

// GetSibling returns the index of the sibling node for the current leaf.
// If the current leaf is on the right, it returns the previous (left) sibling by subtracting 1.
// If the current leaf is on the left, it returns the next (right) sibling by adding 1.
func (n *nodeIndex[TI]) GetSibling() types.Index[TI] {
	shift := types.BitLeftShift[TI](n.GetHeight() + 1)
	types.BitXor(n.value, shift)
	return NodeIndex[TI](types.BitXor(n.value, shift))
}

// RightUp moves the current node to its parent if it's a right child in the tree hierarchy.
func (n *nodeIndex[TI]) RightUp() types.Index[TI] {
	shift := types.BitLeftShift[TI](n.GetHeight() + 1)

	if types.BitAnd[TI](n.value, shift) == shift {
		value := types.BitXor[TI](n.value, types.BitLeftShift[TI](n.GetHeight()))
		if !types.IsNull(value) {
			return NodeIndex[TI](value)
		}
	}
	return nil
}

// Up returns the index of the parent node of the current node in the tree.
func (n *nodeIndex[TI]) Up() types.Index[TI] {
	node := NodeIndex[TI](n.Index())
	if !n.IsRight() {
		node = node.GetSibling()
	}
	return node.RightUp()
}

// IsRight checks if the current node is a right child in the tree hierarchy.
func (n *nodeIndex[TI]) IsRight() bool {
	shift := types.BitLeftShift[TI](n.GetHeight() + 1)
	return types.BitAnd[TI](n.value, shift) == shift
}

// Top returns the index of the "top" ancestor of the current node, climbing the tree structure.
func (n *nodeIndex[TI]) Top() types.Index[TI] {
	shift := types.BitLeftShift[TI](n.GetHeight())
	value := n.value
	result := value
	for !types.IsNull(result) && types.Equal(types.BitAnd(value, shift), shift) {
		result = value
		value = types.BitXor(value, shift)
		shift = types.BitLeft[TI](shift)
	}
	return NodeIndex[TI](result)
}

// Index returns the index value of the current node.
func (n *nodeIndex[TI]) Index() TI {
	return n.value
}

// Children returns the indexes of the children nodes of the current node.
func (n *nodeIndex[TI]) Children() []types.Index[TI] {

	h := n.GetHeight()
	pow := uint64(math.Pow(2, float64(h))) / 2
	index := n.value
	if h == 0 {
		return []types.Index[TI]{LeafIndex(types.SubtractUint64(index, 1)), LeafIndex(index)}
	}
	return []types.Index[TI]{NodeIndex(types.SubtractUint64(index, pow)), NodeIndex(types.SubtractUint64(index, -pow))}
}

// IsLeaf checks if the current node is a leaf node.
func (n *nodeIndex[TI]) IsLeaf() bool {
	return false
}
