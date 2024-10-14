package index_test

import (
	"github.com/dk-open/go-mmr/merkle/index"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getLeafIndex(value uint32) index.Index[uint32] {
	return index.LeafIndex[uint32](value)
}

func TestLeafIndex_Creation(t *testing.T) {
	// Test creating a leaf index
	value := uint32(5)
	leaf := getLeafIndex(value)
	assert.Equal(t, value, leaf.Index(), "Leaf index should be set correctly")
}

func TestLeafIndex_IsLeaf(t *testing.T) {
	// Test if the node is correctly identified as a leaf
	leaf := getLeafIndex(uint32(3))
	assert.True(t, leaf.IsLeaf(), "Leaf index should return true for IsLeaf")
}

func TestLeafIndex_GetHeight(t *testing.T) {
	// Test that the height of a leaf node is always 0
	leaf := getLeafIndex(uint32(3))
	assert.Equal(t, 0, leaf.GetHeight(), "Leaf height should always be 0")
}

func TestLeafIndex_LeftBranch(t *testing.T) {
	// Test left branch calculation
	leaf := getLeafIndex(uint32(2)) // Left child
	leftBranch := leaf.LeftBranch()
	assert.NotNil(t, leftBranch, "Left branch should exist for non-zero, non-right node")
	assert.Equal(t, uint32(1), leftBranch.Index(), "Left branch should be correctly calculated")

	// Test no left branch for node 0
	leaf = getLeafIndex(uint32(0)) // Node 0 has no left branch
	leftBranch = leaf.LeftBranch()
	assert.Nil(t, leftBranch, "Left branch should not exist for node 0")
}

func TestLeafIndex_GetSibling(t *testing.T) {
	// Test sibling calculation for both left and right nodes
	leaf := getLeafIndex(uint32(3)) // Right sibling
	sibling := leaf.GetSibling()
	assert.Equal(t, uint32(2), sibling.Index(), "Right sibling should be calculated correctly")

	leaf = getLeafIndex(uint32(2)) // Left sibling
	sibling = leaf.GetSibling()
	assert.Equal(t, uint32(3), sibling.Index(), "Left sibling should be calculated correctly")
}

func TestLeafIndex_RightUp(t *testing.T) {
	// Test moving to the parent node from a right child
	leaf := getLeafIndex(uint32(3)) // Right child of parent
	parent := leaf.RightUp()
	assert.NotNil(t, parent, "Parent should exist")
	assert.Equal(t, uint32(3), parent.Index(), "Parent node should be calculated correctly")

	// Test when there is no parent
	leaf = getLeafIndex(uint32(0)) // Node 0 has no parent
	parent = leaf.RightUp()
	assert.Nil(t, parent, "There should be no parent for node 0")
}

func TestLeafIndex_Up(t *testing.T) {
	// Test moving up in the hierarchy
	leaf := getLeafIndex(uint32(3)) // Right child
	parent := leaf.Up()
	assert.Equal(t, uint32(3), parent.Index(), "Parent node should be calculated correctly")

	leaf = getLeafIndex(uint32(2)) // Left child
	parent = leaf.Up()
	assert.Equal(t, uint32(3), parent.Index(), "Parent node should be calculated correctly")
}

func TestLeafIndex_IsRight(t *testing.T) {
	// Test checking if a leaf is a right child
	leaf := getLeafIndex(uint32(3)) // Right child
	assert.True(t, leaf.IsRight(), "Leaf should be identified as a right child")

	leaf = getLeafIndex(uint32(2)) // Left child
	assert.False(t, leaf.IsRight(), "Leaf should be identified as a left child")
}

func TestLeafIndex_Top(t *testing.T) {
	// Test finding the top ancestor
	leaf := getLeafIndex(uint32(3)) // Start at leaf 3
	top := leaf.Top()
	assert.Equal(t, uint32(2), top.Index(), "Top ancestor should be calculated correctly")

	leaf = getLeafIndex(uint32(2)) // Start at leaf 2
	top = leaf.Top()
	assert.Equal(t, uint32(2), top.Index(), "Top ancestor should return the leaf itself when it is the top")
}

func TestLeafIndex_Children(t *testing.T) {
	// Test that leaf nodes have no children
	leaf := getLeafIndex(uint32(3))
	children := leaf.Children()
	assert.Nil(t, children, "Leaf nodes should have no children")
}

func TestLeafIndex_Key(t *testing.T) {
	// Test generating the key for the leaf
	leaf := getLeafIndex(uint32(5))
	key := leaf.Key()
	assert.Equal(t, "leaf_5", key, "Leaf key should be generated correctly")
}
