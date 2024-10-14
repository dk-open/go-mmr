package types

import "bytes"

type HashType interface {
	Hash128 | Hash160 | Hash224 | Hash256 | Hash384 | Hash512 |
		uint16 | uint32 | uint64 | int16 | int32 | int64 | int | uint | string
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
	case uint16:
		return []byte{byte(v >> 8), byte(v)}, nil
	case uint32:
		return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
	case uint64:
		return []byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}, nil
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

// BufferWrite writes the given hash or other supported types to the buffer.
func BufferWrite[THash HashType](buf *bytes.Buffer, hash THash) error {
	data, err := HashBytes(hash)
	if err != nil {
		return err
	}

	_, err = buf.Write(data)
	return err
}

// BufferRead reads the given hash or other supported types from the buffer.
func BufferRead[THash HashType](buf *bytes.Reader) (res THash, err error) {
	switch v := any(res).(type) {
	case Hash128:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case Hash160:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case Hash224:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case Hash256:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case Hash384:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case Hash512:
		_, err = buf.Read(v[:])
		return any(v).(THash), err
	case uint16:
		var b [2]byte
		_, err = buf.Read(b[:])
		return any(uint16(b[0])<<8 | uint16(b[1])).(THash), err
	case uint32:
		var b [4]byte
		_, err = buf.Read(b[:])
		return any(uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])).(THash), err
	case uint64:
		var b [8]byte
		_, err = buf.Read(b[:])
		return any(uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
			uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7])).(THash), err
	case int16:
		var b [2]byte
		_, err = buf.Read(b[:])
		return any(int16(b[0])<<8 | int16(b[1])).(THash), err
	case int32:
		var b [4]byte
		_, err = buf.Read(b[:])
		return any(int32(b[0])<<24 | int32(b[1])<<16 | int32(b[2])<<8 | int32(b[3])).(THash), err
	case int64:
		var b [8]byte
		_, err = buf.Read(b[:])
		return any(int64(b[0])<<56 | int64(b[1])<<48 | int64(b[2])<<40 | int64(b[3])<<32 |
			int64(b[4])<<24 | int64(b[5])<<16 | int64(b[6])<<8 | int64(b[7])).(THash), err
	case int:
		var b [4]byte
		_, err = buf.Read(b[:])
		return any(int(b[0])<<24 | int(b[1])<<16 | int(b[2])<<8 | int(b[3])).(THash), err // Assuming 32-bit int
	case uint:
		var b [4]byte
		_, err = buf.Read(b[:])
		return any(uint(b[0])<<24 | uint(b[1])<<16 | uint(b[2])<<8 | uint(b[3])).(THash), err // Assuming 32-bit uint
	case string:
		data := make([]byte, buf.Len())
		_, err = buf.Read(data)
		return any(string(data)).(THash), err
	default:
		return res, ErrTypeMismatch
	}
}
