package types

import (
	"math/big"
)

// NumericValue - Separate type for numeric values (excluding big.Int)
type NumericValue interface {
	int | int32 | int64 | uint | uint32 | uint64
}

// IndexValue - Separate type for index values
type IndexValue interface {
	int | int32 | int64 | uint32 | uint64 | *big.Int
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
	//Key() IK
}

// SubtractUint64 - Subtraction function
func SubtractUint64[IV IndexValue](a IV, b uint64) IV {
	switch v := any(a).(type) {
	case int:
		return any(numericSub(v, int(b))).(IV)
	case int32:
		return any(numericSub(v, int32(b))).(IV)
	case int64:
		return any(numericSub(v, int64(b))).(IV)
	case uint:
		return any(numericSub(v, uint(b))).(IV)
	case uint32:
		return any(numericSub(v, uint32(b))).(IV)
	case uint64:
		return any(numericSub(v, b)).(IV)

	case *big.Int:
		return any(new(big.Int).Sub(v, any(b).(*big.Int))).(IV)
	default:
		panic("unsupported type")
	}
}

// Add - Addition function
func Add[IV IndexValue](a, b IV) IV {
	switch v := any(a).(type) {
	case int:
		return any(numericAdd(v, any(b).(int))).(IV)
	case int32:
		return any(numericAdd(v, any(b).(int32))).(IV)
	case int64:
		return any(numericAdd(v, any(b).(int64))).(IV)
	case uint:
		return any(numericAdd(v, any(b).(uint))).(IV)
	case uint32:
		return any(numericAdd(v, any(b).(uint32))).(IV)
	case uint64:
		return any(numericAdd(v, any(b).(uint64))).(IV)
	case *big.Int:
		return any(new(big.Int).Add(v, any(b).(*big.Int))).(IV)
	default:
		panic("unsupported type")
	}
}

// AddInt  - Addition function
func AddInt[IV IndexValue](a IV, value int) IV {
	switch v := any(a).(type) {
	case int:
		return any(numericAdd(v, value)).(IV)
	case int32:
		return any(numericAdd(v, int32(value))).(IV)
	case int64:
		return any(numericAdd(v, int64(value))).(IV)
	case uint:
		return any(numericAdd(v, uint(value))).(IV)
	case uint32:
		return any(numericAdd(v, uint32(value))).(IV)
	case uint64:
		return any(numericAdd(v, uint64(value))).(IV)
	case *big.Int:
		return any(new(big.Int).Add(v, big.NewInt(int64(value)))).(IV)
	default:
		panic("unsupported type")
	}
}

// Value - Convert int to IndexValue
func Value[IV IndexValue](value uint64) (res IV) {
	switch any(res).(type) {
	case int:
		return any(int(value)).(IV)
	case int32:
		return any(int32(value)).(IV)
	case int64:
		return any(int64(value)).(IV)
	case uint:
		return any(uint(value)).(IV)
	case uint32:
		return any(uint32(value)).(IV)
	case uint64:
		return any(value).(IV)
	case *big.Int:
		return any(big.NewInt(int64(value))).(IV)
	default:
		panic("unsupported type")
	}
}

// GreaterThan function for IndexValue
func GreaterThan[IV IndexValue](a, b IV) bool {
	switch v := any(a).(type) {
	case int:
		return numericGreaterThan(v, any(b).(int))
	case int32:
		return numericGreaterThan(v, any(b).(int32))
	case int64:
		return numericGreaterThan(v, any(b).(int64))
	case uint:
		return numericGreaterThan(v, any(b).(uint))
	case uint32:
		return numericGreaterThan(v, any(b).(uint32))
	case uint64:
		return numericGreaterThan(v, any(b).(uint64))
	case *big.Int:
		return v.Cmp(any(b).(*big.Int)) > 0
	default:
		panic("unsupported type")
	}
}

// IndexUint64 - Convert IndexValue to uint64
func IndexUint64[IV IndexValue](value IV) uint64 {
	switch v := any(value).(type) {
	case int:
		return uint64(v)
	case int32:
		return uint64(v)
	case int64:
		return uint64(v)
	case uint:
		return uint64(v)
	case uint32:
		return uint64(v)
	case uint64:
		return v
	case *big.Int:
		return v.Uint64()
	default:
		panic("unsupported type")
	}

}

// BitLeftShift - Bitwise right shift function
func BitLeftShift[IV IndexValue](size int) (res IV) {
	switch any(res).(type) {
	case int, int32, int64, uint, uint32, uint64:
		return Value[IV](1 << size)
	case *big.Int:
		shifted := new(big.Int).Lsh(big.NewInt(1), uint(size)) // 1 << shiftAmount

		return any(shifted).(IV)
	default:
		panic("unsupported type")
	}
}

// BitLeft - Bitwise left shift function
func BitLeft[IV IndexValue](value IV) (res IV) {
	switch v := any(value).(type) {
	case int:
		return any(v << 1).(IV)
	case int32:
		return any(v << 1).(IV)
	case int64:
		return any(v << 1).(IV)
	case uint:
		return any(v << 1).(IV)
	case uint32:
		return any(v << 1).(IV)
	case uint64:
		return any(v << 1).(IV)
	case *big.Int:
		shifted := new(big.Int).Lsh(any(value).(*big.Int), uint(1)) // 1 << shiftAmount

		return any(shifted).(IV)
	default:
		panic("unsupported type")
	}
}

// IsNull - Check if the value is null
func IsNull[IV IndexValue](value IV) bool {
	switch v := any(value).(type) {
	case int:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case *big.Int:
		return v.Cmp(big.NewInt(0)) == 0
	default:
		panic("unsupported type")
	}
}

// BitAnd - Bitwise AND function
func BitAnd[IV IndexValue](a, b IV) (res IV) {
	switch v := any(a).(type) {
	case int:
		return any(v & any(b).(int)).(IV)
	case int32:
		return any(v & any(b).(int32)).(IV)
	case int64:
		return any(v & any(b).(int64)).(IV)
	case uint:
		return any(v & any(b).(uint)).(IV)
	case uint32:
		return any(v & any(b).(uint32)).(IV)
	case uint64:
		return any(v & any(b).(uint64)).(IV)
	case *big.Int:
		return any(new(big.Int).And(v, any(b).(*big.Int))).(IV)
	default:
		panic("unsupported type")
	}
}

// BitXor - XOR function
func BitXor[IV IndexValue](a, b IV) (res IV) {
	switch v := any(a).(type) {
	case int:
		return any(v ^ any(b).(int)).(IV)
	case int32:
		return any(v ^ any(b).(int32)).(IV)
	case int64:
		return any(v ^ any(b).(int64)).(IV)
	case uint:
		return any(v ^ any(b).(uint)).(IV)
	case uint32:
		return any(v ^ any(b).(uint32)).(IV)
	case uint64:
		return any(v ^ any(b).(uint64)).(IV)
	case *big.Int:
		return any(new(big.Int).Xor(v, any(b).(*big.Int))).(IV)
	default:
		panic("unsupported type")
	}
}

// Equal - Comparison function
func Equal[IV IndexValue](a, b IV) bool {
	switch v := any(a).(type) {
	case int, int32, int64, uint, uint32, uint64:
		return a == b
	case *big.Int:
		return v.Cmp(any(b).(*big.Int)) == 0
	default:
		panic("unsupported type")
	}
}

// Generic comparison function for numeric types
func numericGreaterThan[IV NumericValue](a, b IV) bool {
	return a > b
}

// Generic subtraction for numeric types
func numericSub[IV NumericValue](a, b IV) IV {
	return a - b
}

// Generic addition for numeric types
func numericAdd[IV NumericValue](a IV, b IV) IV {
	return a + b
}
