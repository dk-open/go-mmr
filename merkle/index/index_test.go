package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHeightMultiType(t *testing.T) {
	assert.Equal(t, 2, getHeight[uint](4))
	assert.Equal(t, 2, getHeight[uint32](4))
	assert.Equal(t, 4, getHeight[int32](16))
	assert.Equal(t, 5, getHeight[uint64](32))
	assert.Equal(t, 1, getHeight[uint16](78))
}

func TestLeaf(t *testing.T) {
	l0 := LeafIndex[int](0)
	assert.Equal(t, false, l0.IsRight())
	assert.Equal(t, 0, l0.Index())
	assert.Equal(t, 1, l0.GetSibling().Index())
	assert.Equal(t, true, l0.GetSibling().IsRight())
	assert.Equal(t, false, l0.Up().IsLeaf())
	assert.Equal(t, 1, l0.Up().Index())
	assert.Equal(t, false, l0.Up().IsRight())
	assert.Equal(t, 1, l0.GetSibling().Top().Index())

	l1 := LeafIndex[int](1)
	assert.Equal(t, true, l1.IsRight())
	assert.Equal(t, false, l1.Up().IsRight())
	assert.Equal(t, false, l1.Up().Up().IsLeaf())
	assert.Equal(t, 2, l1.Up().Up().Index())

	l3 := LeafIndex[int](3)
	assert.Equal(t, 2, l3.Top().Index())

}

func TestPeaks(t *testing.T) {
	l0 := LeafIndex[int](0)
	peaks0 := GetPeaks(l0)
	assert.Equal(t, 1, len(peaks0))

	l1 := LeafIndex[int](1)
	peaks1 := GetPeaks(l1)
	assert.Equal(t, 1, len(peaks1))

	l2 := LeafIndex[int](2)
	peaks2 := GetPeaks(l2)
	assert.Equal(t, 2, len(peaks2))
	assert.Equal(t, 2, peaks2[0].Index())
	assert.Equal(t, 1, peaks2[1].Index())
	l3 := LeafIndex[int](3)
	peaks3 := GetPeaks(l3)
	assert.Equal(t, 1, len(peaks3))
	assert.Equal(t, 2, peaks3[0].Index())

	p5 := LeafIndex[int](5)
	peaks5 := GetPeaks(p5)
	assert.Equal(t, 2, len(peaks5))
	assert.Equal(t, 5, peaks5[0].Index())
	assert.Equal(t, 2, peaks5[1].Index())

	p6 := LeafIndex[int](6)
	peaks6 := GetPeaks(p6)
	assert.Equal(t, 3, len(peaks6))
	assert.Equal(t, 6, peaks6[0].Index())
	assert.Equal(t, 5, peaks6[1].Index())
	assert.Equal(t, 2, peaks6[2].Index())

	p7 := LeafIndex[int](7)
	peaks7 := GetPeaks(p7)
	assert.Equal(t, 1, len(peaks7))
	assert.Equal(t, 4, peaks7[0].Index())

	p8 := LeafIndex[int](8)
	peaks8 := GetPeaks(p8)
	assert.Equal(t, 2, len(peaks8))
	assert.Equal(t, 8, peaks8[0].Index())
	assert.Equal(t, 4, peaks8[1].Index())

	p9 := LeafIndex[int](9)
	peaks9 := GetPeaks(p9)
	assert.Equal(t, 2, len(peaks9))
	assert.Equal(t, 9, peaks9[0].Index())
	assert.Equal(t, 4, peaks9[1].Index())

	p10 := LeafIndex[int](10)
	peaks10 := GetPeaks(p10)
	assert.Equal(t, 3, len(peaks10))
	assert.Equal(t, 10, peaks10[0].Index())
	assert.Equal(t, 9, peaks10[1].Index())
	assert.Equal(t, 4, peaks10[2].Index())

	p14 := LeafIndex[int](14)
	peaks14 := GetPeaks(p14)
	assert.Equal(t, 4, len(peaks14))
	assert.Equal(t, 14, peaks14[0].Index())
	assert.Equal(t, 13, peaks14[1].Index())
	assert.Equal(t, 10, peaks14[2].Index())
	assert.Equal(t, 4, peaks14[3].Index())

}
