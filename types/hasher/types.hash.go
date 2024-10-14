package hasher

import (
	"crypto/sha256"
	"crypto/sha512"
	"github.com/dk-open/go-mmr/types"
	"github.com/minio/blake2b-simd"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	"hash"
)

// Sha256 creates a Hasher for SHA-256.
func Sha256(values ...[]byte) types.Hash256 {
	return types.Hash256(sumHashes(sha256.New(), values...))
}

func Sha512(values ...[]byte) types.Hash512 {
	return types.Hash512(sumHashes(sha512.New(), values...))
}

func Sha3_256(values ...[]byte) types.Hash256 {
	return types.Hash256(sumHashes(sha3.New256(), values...))
}

func Sha3_384(values ...[]byte) types.Hash384 {
	return types.Hash384(sumHashes(sha3.New384(), values...))
}

func Sha3_512(values ...[]byte) types.Hash512 {
	return types.Hash512(sumHashes(sha3.New512(), values...))
}

func Blake2b_256(values ...[]byte) types.Hash256 {
	return types.Hash256(sumHashes(blake2b.New256(), values...))
}

func Blake2b_512(values ...[]byte) types.Hash512 {
	return types.Hash512(sumHashes(blake2b.New512(), values...))
}

func Ripemd160(values ...[]byte) types.Hash160 {
	return types.Hash160(sumHashes(ripemd160.New(), values...))
}

// Argon2 creates a Hasher that uses the Argon2id variant.
func Argon2(values ...[]byte) types.Hash256 {
	salt := make([]byte, 16) // Normally, you would generate a random salt here
	hash := argon2.IDKey(values[0], salt, 1, 64*1024, 4, 32)
	return types.Hash256(hash)
}

// Blake3 creates a Hasher that uses BLAKE3.
func Blake3(values ...[]byte) types.Hash256 {
	return types.Hash256(sumHashes(blake3.New(), values...))
}

func sumHashes(h hash.Hash, values ...[]byte) []byte {
	for i := len(values) - 1; i >= 0; i-- {
		h.Write(values[i])
	}
	return h.Sum([]byte{})
}
