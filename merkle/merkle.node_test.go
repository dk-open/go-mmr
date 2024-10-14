package merkle_test

import (
	"bytes"
	"fmt"
	"github.com/dk-open/go-mmr/types/hasher"
	"testing"

	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/types"
	"github.com/stretchr/testify/assert"
)

func TestLeafNodeEncoding(t *testing.T) {
	t.Run("Test LeafNode Encoding for Index 0", func(t *testing.T) {
		ln0 := merkle.Node[types.Hash224]()
		assert.NotNil(t, ln0, "LeafNode should not be nil for index 0")

		data, err := ln0.MarshalBinary()
		assert.NoError(t, err, "Marshaling LeafNode failed")

		ln0_u, err := merkle.NodeFromBinary[types.Hash224](data)
		assert.NoError(t, err, "Unmarshaling LeafNode failed")
		assert.NotNil(t, ln0_u, "Unmarshaled node should not be nil")

		assert.Equal(t, ln0.Children(), ln0_u.Children(), "Children mismatch after unmarshaling")
	})

	t.Run("Test LeafNode Encoding for Index 1", func(t *testing.T) {
		ln1 := merkle.Node[types.Hash224]()
		assert.NotNil(t, ln1, "LeafNode should not be nil for index 1")

		data, err := ln1.MarshalBinary()
		fmt.Printf("Marshalled %x\n", data)
		assert.NoError(t, err, "Marshaling LeafNode failed")

		ln1_u, err := merkle.NodeFromBinary[types.Hash224](data)
		assert.NoError(t, err, "Unmarshaling LeafNode failed")
		assert.NotNil(t, ln1_u, "Unmarshaled node should not be nil")

		assert.Equal(t, ln1.Children(), ln1_u.Children(), "Children mismatch after unmarshaling")
	})
}

func TestLeafNodeWithMultipleChildren(t *testing.T) {
	t.Run("Test LeafNode with Multiple Children", func(t *testing.T) {
		left := hasher.Blake2b_256([]byte("test left"))
		right := hasher.Blake2b_256([]byte("test right"))

		ln := merkle.Node[types.Hash256](left, right)

		assert.Equal(t, 2, len(ln.Children()), "LeafNode should have 2 children")

		data, err := ln.MarshalBinary()
		assert.NoError(t, err, "Marshaling LeafNode with multiple children failed")

		lnU, err := merkle.NodeFromBinary[types.Hash256](data)
		//leftU, ok := lnU.Left()
		//assert.True(t, ok, "Left child not found after unmarshaling")

		//rightU, ok := lnU.Right()
		//assert.True(t, ok, "Right child not found after unmarshaling")
		assert.NoError(t, err, "Unmarshaling LeafNode with multiple children failed")
		assert.Equal(t, 2, len(lnU.Children()), "Children mismatch after unmarshaling")
		assert.Equal(t, left, lnU.Children()[0], "Left child mismatch after unmarshaling")
		//assert.Equal(t, left, leftU, "Left child mismatch after unmarshaling")
		assert.Equal(t, right, lnU.Children()[1], "Right child mismatch after unmarshaling")
		//assert.Equal(t, right, rightU, "Right child mismatch after unmarshaling")
	})
}

func TestMarshalingErrors(t *testing.T) {
	t.Run("Test Invalid Node Binary Data", func(t *testing.T) {
		invalidData := []byte{0x01, 0x02, 0x03} // Some invalid data
		_, err := merkle.NodeFromBinary[types.Hash224](invalidData)
		assert.Error(t, err, "Expected error when unmarshaling invalid data")
	})
}

func TestNode_Children(t *testing.T) {
	// Create a node with mock children
	left := types.Hash256{1, 2, 3}
	right := types.Hash256{4, 5, 6}

	n := merkle.Node[types.Hash256](left, right)

	// Test that Children returns the correct values
	children := n.Children()
	assert.Equal(t, left, children[0], "Left child should be correct")
	assert.Equal(t, right, children[1], "Right child should be correct")
}

func TestNode_SetLeftAndRight(t *testing.T) {
	// Create a node
	left := types.Hash256{1, 2, 3}
	right := types.Hash256{4, 5, 6}

	n := merkle.Node(left, right)

	// Set new left and right children
	newLeft := types.Hash256{7, 8, 9}
	newRight := types.Hash256{10, 11, 12}

	n.SetLeft(newLeft)
	n.SetRight(newRight)

	// Test that the children have been updated
	children := n.Children()
	assert.Equal(t, newLeft, children[0], "Left child should be updated")
	assert.Equal(t, newRight, children[1], "Right child should be updated")
}

func TestNode_MarshalBinary(t *testing.T) {
	// Create a node
	left := types.Hash256{1, 2, 3}
	right := types.Hash256{4, 5, 6}

	n := merkle.Node(left, right)

	// Marshal to binary
	data, err := n.MarshalBinary()
	assert.NoError(t, err, "Marshaling should not return an error")

	// Verify the length of the marshaled data
	expectedLength := len(left) + len(right)
	assert.Len(t, data, expectedLength, "Marshaled binary data length should be correct")
}

func TestNode_UnmarshalBinary(t *testing.T) {
	// Create mock binary data for two children
	left := types.Hash256{1, 2, 3}
	right := types.Hash256{4, 5, 6}

	var buf bytes.Buffer
	types.BufferWrite(&buf, left)
	types.BufferWrite(&buf, right)

	// Unmarshal the binary data into a node
	n, err := merkle.NodeFromBinary[types.Hash256](buf.Bytes())
	assert.NoError(t, err, "Unmarshaling should not return an error")

	// Test that the children have been correctly unmarshaled
	children := n.Children()
	assert.Equal(t, left, children[0], "Left child should be correctly unmarshaled")
	assert.Equal(t, right, children[1], "Right child should be correctly unmarshaled")
}
