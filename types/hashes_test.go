package types_test

import (
	"bytes"
	"fmt"
	"github.com/dk-open/go-mmr/types"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var tests = []struct {
	input    interface{}
	expected []byte
}{
	{types.Hash128{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}},                                                                                                                                                                                                                                                                                                                                                                                                                                 // 12 bytes
	{types.Hash160{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}},                                                                                                                                                                                                                                                                                                                                                                 // 20 bytes
	{types.Hash224{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28}},                                                                                                                                                                                                                                                                                                 // 28 bytes
	{types.Hash256{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}},                                                                                                                                                                                                                                                                 // 32 bytes
	{types.Hash384{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48}},                                                                                                                                 // 48 bytes
	{types.Hash512{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64}}, // 64 bytes
	// Primitive types
	{uint16(258), []byte{1, 2}},
	{uint32(16909060), []byte{1, 2, 3, 4}},
	{uint64(72623859790382856), []byte{1, 2, 3, 4, 5, 6, 7, 8}},
	{int16(258), []byte{1, 2}},
	{int32(16909060), []byte{1, 2, 3, 4}},
	{int64(72623859790382856), []byte{1, 2, 3, 4, 5, 6, 7, 8}},
	{int(16909060), []byte{1, 2, 3, 4}},  // Assuming 32-bit int
	{uint(16909060), []byte{1, 2, 3, 4}}, // Assuming 32-bit uint
	{string("test"), []byte("test")},
}

func TestHashBytes(t *testing.T) {
	for _, test := range tests {
		var result []byte
		var err error
		switch test.input.(type) {
		case types.Hash128:
			result, err = types.HashBytes[types.Hash128](test.input.(types.Hash128))
		case types.Hash160:
			result, err = types.HashBytes[types.Hash160](test.input.(types.Hash160))
		case types.Hash224:
			result, err = types.HashBytes[types.Hash224](test.input.(types.Hash224))
		case types.Hash256:
			result, err = types.HashBytes[types.Hash256](test.input.(types.Hash256))
		case types.Hash384:
			result, err = types.HashBytes[types.Hash384](test.input.(types.Hash384))
		case types.Hash512:
			result, err = types.HashBytes[types.Hash512](test.input.(types.Hash512))
		case uint16:
			result, err = types.HashBytes[uint16](test.input.(uint16))
		case uint32:
			result, err = types.HashBytes[uint32](test.input.(uint32))
		case uint64:
			result, err = types.HashBytes[uint64](test.input.(uint64))
		case int16:
			result, err = types.HashBytes[int16](test.input.(int16))
		case int32:
			result, err = types.HashBytes[int32](test.input.(int32))
		case int64:
			result, err = types.HashBytes[int64](test.input.(int64))
		case int:
			result, err = types.HashBytes[int](test.input.(int))
		case uint:
			result, err = types.HashBytes[uint](test.input.(uint))
		case string:
			result, err = types.HashBytes[string](test.input.(string))

		}
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		fmt.Println("res", result, test.input)

		if !assert.Equal(t, test.expected, result) {
			t.Errorf("HashBytes(%v) = %v; want %v", test.input, result, test.expected)
		}
	}
}

func TestBufferWrite(t *testing.T) {
	for _, test := range tests {
		var buf bytes.Buffer
		var err error
		switch test.input.(type) {
		case types.Hash128:
			err = types.BufferWrite[types.Hash128](&buf, test.input.(types.Hash128))
		case types.Hash160:
			err = types.BufferWrite[types.Hash160](&buf, test.input.(types.Hash160))
		case types.Hash224:
			err = types.BufferWrite[types.Hash224](&buf, test.input.(types.Hash224))
		case types.Hash256:
			err = types.BufferWrite[types.Hash256](&buf, test.input.(types.Hash256))
		case types.Hash384:
			err = types.BufferWrite[types.Hash384](&buf, test.input.(types.Hash384))
		case types.Hash512:
			err = types.BufferWrite[types.Hash512](&buf, test.input.(types.Hash512))
		case uint16:
			err = types.BufferWrite[uint16](&buf, test.input.(uint16))
		case uint32:
			err = types.BufferWrite[uint32](&buf, test.input.(uint32))
		case uint64:
			err = types.BufferWrite[uint64](&buf, test.input.(uint64))
		case int16:
			err = types.BufferWrite[int16](&buf, test.input.(int16))
		case int32:
			err = types.BufferWrite[int32](&buf, test.input.(int32))
		case int64:
			err = types.BufferWrite[int64](&buf, test.input.(int64))
		case int:
			err = types.BufferWrite[int](&buf, test.input.(int))
		case uint:
			err = types.BufferWrite[uint](&buf, test.input.(uint))
		case string:
			err = types.BufferWrite[string](&buf, test.input.(string))
		}
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if !assert.Equal(t, test.expected, buf.Bytes()) {
			t.Errorf("HashBytes(%v) = %v; want %v", test.input, buf.Bytes(), test.expected)
		}
	}
}

func TestBufferRead(t *testing.T) {
	for _, test := range tests {
		buf := bytes.NewReader(test.expected)
		var result interface{}
		var err error

		switch test.input.(type) {
		case types.Hash128:
			result, err = types.BufferRead[types.Hash128](buf)
		case types.Hash160:
			result, err = types.BufferRead[types.Hash160](buf)
		case types.Hash224:
			result, err = types.BufferRead[types.Hash224](buf)
		case types.Hash256:
			result, err = types.BufferRead[types.Hash256](buf)
		case types.Hash384:
			result, err = types.BufferRead[types.Hash384](buf)
		case types.Hash512:
			result, err = types.BufferRead[types.Hash512](buf)
		case uint16:
			result, err = types.BufferRead[uint16](buf)
		case uint32:
			result, err = types.BufferRead[uint32](buf)
		case uint64:
			result, err = types.BufferRead[uint64](buf)
		case int16:
			result, err = types.BufferRead[int16](buf)
		case int32:
			result, err = types.BufferRead[int32](buf)
		case int64:
			result, err = types.BufferRead[int64](buf)
		case int:
			result, err = types.BufferRead[int](buf)
		case uint:
			result, err = types.BufferRead[uint](buf)
		case string:
			result, err = types.BufferRead[string](buf)
		default:
			t.Errorf("Unsupported type: %T", test.expected)
		}

		// Error checking
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Equality check
		if !reflect.DeepEqual(result, test.input) {
			t.Errorf("BufferRead(%v) = %v; want %v", test.expected, result, test.input)
		}
	}
}
