package merkle_test

import (
	"github.com/dk-open/go-mmr/types/hasher"
	"testing"

	"github.com/dk-open/go-mmr/merkle"
	"github.com/dk-open/go-mmr/types"
	"github.com/stretchr/testify/assert"
)

func TestLeafNodeEncoding(t *testing.T) {
	t.Run("Test LeafNode Encoding for Index 0", func(t *testing.T) {
		ln0 := merkle.LeafNode[uint32, types.Hash224](0)
		assert.NotNil(t, ln0, "LeafNode should not be nil for index 0")
		assert.Zero(t, ln0.Index(), "LeafNode index should be 0")
		assert.Zero(t, ln0.Children(), "LeafNode should have no children")
		assert.True(t, ln0.IsLeaf(), "LeafNode should be marked as a leaf")

		data, err := ln0.MarshalBinary()
		assert.NoError(t, err, "Marshaling LeafNode failed")
		assert.NotNil(t, data, "Marshaled data should not be nil")

		ln0_u, err := merkle.NodeFromBinary[uint32, types.Hash224](data)
		assert.NoError(t, err, "Unmarshaling LeafNode failed")
		assert.NotNil(t, ln0_u, "Unmarshaled node should not be nil")

		assert.Equal(t, ln0.Index(), ln0_u.Index(), "Index mismatch after unmarshaling")
		assert.Equal(t, ln0.Children(), ln0_u.Children(), "Children mismatch after unmarshaling")
		assert.Equal(t, ln0.IsLeaf(), ln0_u.IsLeaf(), "IsLeaf mismatch after unmarshaling")
	})

	t.Run("Test LeafNode Encoding for Index 1", func(t *testing.T) {
		ln1 := merkle.LeafNode[uint32, types.Hash224](1)
		assert.NotNil(t, ln1, "LeafNode should not be nil for index 1")
		assert.Equal(t, uint32(1), ln1.Index(), "LeafNode index should be 1")
		assert.Zero(t, ln1.Children(), "LeafNode should have no children")
		assert.True(t, ln1.IsLeaf(), "LeafNode should be marked as a leaf")

		data, err := ln1.MarshalBinary()
		assert.NoError(t, err, "Marshaling LeafNode failed")
		assert.NotNil(t, data, "Marshaled data should not be nil")

		ln1_u, err := merkle.NodeFromBinary[uint32, types.Hash224](data)
		assert.NoError(t, err, "Unmarshaling LeafNode failed")
		assert.NotNil(t, ln1_u, "Unmarshaled node should not be nil")

		assert.Equal(t, ln1.Index(), ln1_u.Index(), "Index mismatch after unmarshaling")
		assert.Equal(t, ln1.Children(), ln1_u.Children(), "Children mismatch after unmarshaling")
		assert.Equal(t, ln1.IsLeaf(), ln1_u.IsLeaf(), "IsLeaf mismatch after unmarshaling")
	})
}

func TestLeafNodeWithMultipleChildren(t *testing.T) {
	t.Run("Test LeafNode with Multiple Children", func(t *testing.T) {
		left := hasher.Blake2b_256([]byte("test left"))
		right := hasher.Blake2b_256([]byte("test right"))

		ln := merkle.Node[uint32, types.Hash256](0, left, right)

		assert.Equal(t, 2, len(ln.Children()), "LeafNode should have 2 children")

		data, err := ln.MarshalBinary()
		assert.NoError(t, err, "Marshaling LeafNode with multiple children failed")

		lnU, err := merkle.NodeFromBinary[uint32, types.Hash256](data)
		leftU, ok := lnU.Left()
		assert.True(t, ok, "Left child not found after unmarshaling")

		rightU, ok := lnU.Right()
		assert.True(t, ok, "Right child not found after unmarshaling")
		assert.NoError(t, err, "Unmarshaling LeafNode with multiple children failed")
		assert.Equal(t, 2, len(lnU.Children()), "Children mismatch after unmarshaling")
		assert.Equal(t, left, lnU.Children()[0], "Left child mismatch after unmarshaling")
		assert.Equal(t, left, leftU, "Left child mismatch after unmarshaling")
		assert.Equal(t, right, lnU.Children()[1], "Right child mismatch after unmarshaling")
		assert.Equal(t, right, rightU, "Right child mismatch after unmarshaling")
	})
}

func TestMarshalingErrors(t *testing.T) {
	t.Run("Test Invalid Node Binary Data", func(t *testing.T) {
		invalidData := []byte{0x01, 0x02, 0x03} // Some invalid data
		_, err := merkle.NodeFromBinary[uint32, types.Hash224](invalidData)
		assert.Error(t, err, "Expected error when unmarshaling invalid data")
	})
}
