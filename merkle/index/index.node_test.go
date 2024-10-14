package index_test

import (
	"github.com/dk-open/go-mmr/merkle/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getNodeIndex(value uint32) index.Index[uint32] {
	return index.NodeIndex[uint32](value)
}

func TestNodeIndex_Creation(t *testing.T) {
	// Test creating a node index
	value := uint32(4)
	node := getNodeIndex(value)
	assert.Equal(t, value, node.Index(), "Node index should be set correctly")
	assert.Equal(t, 2, node.GetHeight(), "Height should be calculated correctly")
}

func TestNodeIndex_LeftBranch(t *testing.T) {
	// Test left branch calculation
	node := getNodeIndex(uint32(6))
	leftBranch := node.LeftBranch()
	assert.NotNil(t, leftBranch, "Left branch should exist")
	assert.Equal(t, uint32(2), leftBranch.Index(), "Left branch should be calculated correctly")
}

func TestNodeIndex_GetSibling(t *testing.T) {
	// Test sibling calculation for both left and right nodes
	node := getNodeIndex(uint32(3)) // Right sibling
	sibling := node.GetSibling()
	assert.Equal(t, uint32(1), sibling.Index(), "Right sibling should be calculated correctly")

	node = getNodeIndex(uint32(2)) // Left sibling
	sibling = node.GetSibling()
	assert.Equal(t, uint32(6), sibling.Index(), "Left sibling should be calculated correctly")
}

func TestNodeIndex_RightUp(t *testing.T) {
	// Test moving to the parent node from a right child
	node := getNodeIndex(uint32(6)) // Right child of parent
	parent := node.RightUp()
	assert.NotNil(t, parent, "Parent should exist")
	assert.Equal(t, uint32(4), parent.Index(), "Parent node should be calculated correctly")

	// Test when already at the top
	node = getNodeIndex(uint32(4)) // Top node
	parent = node.RightUp()
	assert.Nil(t, parent, "There should be no parent when already at the top")
}

func TestNodeIndex_Up(t *testing.T) {
	// Test moving up in the hierarchy
	node := getNodeIndex(uint32(3)) // Right child
	parent := node.Up()
	assert.Equal(t, uint32(2), parent.Index(), "Parent node should be calculated correctly")

	node = getNodeIndex(uint32(2)) // Left child
	parent = node.Up()
	assert.Equal(t, uint32(4), parent.Index(), "Parent node should be calculated correctly")
}

func TestNodeIndex_IsRight(t *testing.T) {
	// Test checking if a node is a right child
	node := getNodeIndex(uint32(3)) // Right child
	assert.True(t, node.IsRight(), "Node should be identified as a right child")

	node = getNodeIndex(uint32(2)) // Left child
	assert.False(t, node.IsRight(), "Node should be identified as a left child")
}

func TestNodeIndex_Top(t *testing.T) {
	// Test finding the top ancestor
	node := getNodeIndex(uint32(6)) // Start at node 6
	top := node.Top()
	assert.Equal(t, uint32(4), top.Index(), "Top ancestor should be calculated correctly")
}

func TestNodeIndex_Children(t *testing.T) {
	// Test retrieving the children of a node
	node := getNodeIndex(uint32(6)) // Node with height > 0
	children := node.Children()
	assert.Len(t, children, 2, "Node should have two children")
	assert.Equal(t, uint32(5), children[0].Index(), "Left child should be calculated correctly")
	assert.Equal(t, uint32(7), children[1].Index(), "Right child should be calculated correctly")

	node = getNodeIndex(uint32(1)) // Node with height 0
	children = node.Children()
	assert.Len(t, children, 2, "Leaf node should have two children")
	assert.Equal(t, uint32(0), children[0].Index(), "Left child should be calculated correctly for leaf")
	assert.Equal(t, uint32(1), children[1].Index(), "Right child should be calculated correctly for leaf")
}

func TestNodeIndex_IsLeaf(t *testing.T) {
	// Test checking if a node is a leaf
	node := getNodeIndex(uint32(1)) // Node at height 0
	assert.False(t, node.IsLeaf(), "Node should not be identified as a leaf")
}

func TestNodeIndex_Key(t *testing.T) {
	// Test generating the key for the node
	node := getNodeIndex(uint32(5))
	key := node.Key()
	assert.Equal(t, "node_5", key, "Node key should be generated correctly")
}
