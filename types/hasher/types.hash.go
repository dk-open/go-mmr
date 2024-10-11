package hasher

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"github.com/dk-open/go-mmr/types"
	"github.com/minio/blake2b-simd"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/salsa20"
	"golang.org/x/crypto/sha3"
	"hash"
	"log"
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

// Salsa20 creates a Hasher using the Salsa20 stream cipher.
func Salsa20(values ...[]byte) types.Hash256 {
	// Salsa20 requires a 32-byte key and 8-byte nonce
	key := [32]byte{}
	nonce := make([]byte, 8)

	// Generate random key and nonce (in practice, securely share or manage these)
	if _, err := rand.Read(key[:]); err != nil {
		log.Fatal(err)
	}
	if _, err := rand.Read(nonce); err != nil {
		log.Fatal(err)
	}

	// Prepare the buffer for the output (same length as input)
	output := make([]byte, len(values[0]))

	// Encrypt or decrypt using salsa20 XOR
	salsa20.XORKeyStream(output, values[0], nonce, &key)

	// Return the result as Hash256 (you can adapt based on what Hash256 is)
	return types.Hash256(output)
}

func sumHashes(h hash.Hash, values ...[]byte) []byte {
	for i := len(values) - 1; i >= 0; i-- {
		h.Write(values[i])
	}
	return h.Sum([]byte{})
}
