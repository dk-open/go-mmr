package index

import (
	"fmt"
)

type nodeIndex[TI Value] struct {
	value  TI
	height int
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
func NodeIndex[TI Value](value TI) Index[TI] {
	res := &nodeIndex[TI]{
		value:  value,
		height: getHeight(value),
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
func (n *nodeIndex[TI]) LeftBranch() Index[TI] {
	distance := getDistance[TI](n.height + 1)
	if n.value > distance {
		return NodeIndex[TI](n.value - distance)
	}
	return nil
}

// GetSibling returns the index of the sibling node for the current leaf.
// If the current leaf is on the right, it returns the previous (left) sibling by subtracting 1.
// If the current leaf is on the left, it returns the next (right) sibling by adding 1.
func (n *nodeIndex[TI]) GetSibling() Index[TI] {
	distance := getDistance[TI](n.height + 1)
	return NodeIndex[TI](n.value ^ distance)
}

// RightUp moves the current node to its parent if it's a right child in the tree hierarchy.
func (n *nodeIndex[TI]) RightUp() Index[TI] {
	distance := getDistance[TI](n.height + 1)
	if n.value&distance == distance {
		if value := n.value ^ (1 << n.height); value > 0 {
			return NodeIndex[TI](value)
		}
	}
	return nil
}

// Up returns the index of the parent node of the current node in the tree.
func (n *nodeIndex[TI]) Up() Index[TI] {
	node := NodeIndex[TI](n.Index())
	if !n.IsRight() {
		node = node.GetSibling()
	}
	return node.RightUp()
}

// IsRight checks if the current node is a right child in the tree hierarchy.
func (n *nodeIndex[TI]) IsRight() bool {
	distance := getDistance[TI](n.height + 1)
	return n.value&distance == distance
}

// Top returns the index of the "top" ancestor of the current node, climbing the tree structure.
func (n *nodeIndex[TI]) Top() Index[TI] {
	shift := TI(1 << n.height)
	value := n.value
	top := value
	for top != 0 && value&shift == shift {
		top = value
		value ^= shift
		shift <<= 1
	}
	return NodeIndex[TI](top)
}

// Index returns the index value of the current node.
func (n *nodeIndex[TI]) Index() TI {
	return n.value
}

// Children returns the indexes of the children nodes of the current node.
func (n *nodeIndex[TI]) Children() []Index[TI] {
	if n.height == 0 {
		return []Index[TI]{LeafIndex(n.value - 1), LeafIndex(n.value)}
	}
	childDistance := getDistance[TI](n.height - 1)
	return []Index[TI]{NodeIndex(n.value - childDistance), NodeIndex(n.value + childDistance)}
}

// IsLeaf checks if the current node is a leaf node.
func (n *nodeIndex[TI]) IsLeaf() bool {
	return false
}

func (n *nodeIndex[TI]) Key() string {
	return fmt.Sprintf("node_%d", n.value)
}

func getDistance[TI Value](height int) TI {
	return TI(1 << height)
}
