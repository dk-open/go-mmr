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
