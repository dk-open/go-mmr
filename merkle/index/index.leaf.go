package index

import (
	"fmt"
)

type leafIndex[TI IndexValue] struct {
	value TI
	//keyF  KeyFunc[TK, TV]
}

// LeafIndex creates a new leaf index with the given value.
func LeafIndex[TI IndexValue](value TI) Index[TI] {
	res := &leafIndex[TI]{
		value: value,
	}
	return res
}

// IsLeaf returns true if the current node is a leaf node in the tree.
func (l *leafIndex[TI]) IsLeaf() bool {
	return true
}

// GetHeight returns the height of the current node in the tree - 0 for leaf nodes.
func (l *leafIndex[TI]) GetHeight() int {
	return 0
}

// LeftBranch calculates the index of the left child (left branch) of the current node.
// It checks if the left child exists based on the node's height and position,
// and returns the left child's index or nil if there is no left branch.
func (l *leafIndex[TI]) LeftBranch() Index[TI] {
	if !l.IsRight() && l.value != 0 {
		return NodeIndex[TI](l.value - 1)
	}
	return nil
}

// GetSibling returns the index of the sibling node for the current leaf.
// If the current leaf is on the right, it returns the previous (left) sibling by subtracting 1.
// If the current leaf is on the left, it returns the next (right) sibling by adding 1.
func (l *leafIndex[TI]) GetSibling() Index[TI] {
	if l.IsRight() {
		return LeafIndex(l.value - 1)
	}
	return LeafIndex(l.value + 1)
}

// RightUp moves the current node to its parent if it's a right child in the tree hierarchy.
func (l *leafIndex[TI]) RightUp() Index[TI] {
	//	value := x.Index()
	if l.IsRight() {
		return NodeIndex(l.value)
	}
	return nil
}

// Up returns the index of the parent node of the current node in the tree.
func (l *leafIndex[TI]) Up() (res Index[TI]) {
	res = l
	if !l.IsRight() {
		res = res.GetSibling()
	}
	return res.RightUp()
}

func (l *leafIndex[TI]) IsRight() bool {
	return l.value&1 == 1
}

// Top returns the index of the "top" ancestor of the current node, climbing the tree structure.
func (l *leafIndex[TI]) Top() Index[TI] {
	if !l.IsRight() {
		return l
	}
	return NodeIndex(l.value).Top()
}

// Index returns the index value of the current node.
func (l *leafIndex[TI]) Index() TI {
	return l.value
}

// Children returns the children of the current node.
func (l *leafIndex[TI]) Children() []Index[TI] {
	return nil //No children for object index
}

func (l *leafIndex[TI]) Key() string {
	return fmt.Sprintf("leaf_%d", l.value)
}
