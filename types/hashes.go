package types

type HashType interface {
	Hash128 | Hash160 | Hash224 | Hash256 | Hash384 | Hash512 |
		uint8 | uint16 | uint32 | uint64 | int8 | int16 | int32 | int64 | int | uint | string
}

type Hash128 [12]byte
type Hash160 [20]byte
type Hash224 [28]byte
type Hash256 [32]byte
type Hash384 [48]byte
type Hash512 [64]byte

type IHash[Key HashType] interface {
	Hash() Key
}

type Hasher[Key HashType] func(values ...[]byte) Key

func HashBytes[TH HashType](value TH) ([]byte, error) {
	switch v := any(value).(type) {
	case Hash128:
		return v[:], nil
	case Hash160:
		return v[:], nil
	case Hash224:
		return v[:], nil
	case Hash256:
		return v[:], nil
	case Hash384:
		return v[:], nil
	case Hash512:
		return v[:], nil
	case uint8:
		return []byte{byte(v)}, nil
	case uint16:
		return []byte{byte(v >> 8), byte(v)}, nil
	case uint32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case uint64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int8:
		return []byte{byte(v)}, nil
	case int16:
		return []byte{byte(v >> 8), byte(v)}, nil
	case int32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case int:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case uint:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case string:
		return []byte(v), nil
	default:
		return nil, ErrTypeMismatch
	}
}

//
//// ReadHash Helper to read the hash field from the buffer (same as before)
//func ReadHash[THash HashType](buf *bytes.Reader) (THash, error) {
//	switch any(th).(type) {
//	case types.Hash160:
//		var left, right types.Hash160
//		copy(left[:], data[:20])
//		copy(right[:], data[20:40])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	case types.Hash256:
//		var left, right types.Hash256
//		copy(left[:], data[:32])
//		copy(right[:], data[32:64])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	case types.Hash512:
//		var left, right types.Hash512
//		copy(left[:], data[:64])
//		copy(right[:], data[64:128])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	case types.Hash128:
//		var left, right types.Hash128
//		copy(left[:], data[:16])
//		copy(right[:], data[16:32])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	case types.Hash224:
//		var left, right types.Hash224
//		copy(left[:], data[:28])
//		copy(right[:], data[28:56])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	case types.Hash384:
//		var left, right types.Hash384
//		copy(left[:], data[:48])
//		copy(right[:], data[48:96])
//		res.Child = []TH{any(left).(TH), any(right).(TH)}
//	default:
//		return res, types.ErrTypeMismatch
//	}
//}
