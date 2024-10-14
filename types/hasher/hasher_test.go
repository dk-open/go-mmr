package hasher_test

import (
	"github.com/dk-open/go-mmr/types/hasher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSha256(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Sha256(value)
	actual := hasher.Sha256(value)
	assert.Equal(t, expected, actual, "SHA-256 hash should be equal")
}

func TestSha512(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Sha512(value)
	actual := hasher.Sha512(value)
	assert.Equal(t, expected, actual, "SHA-512 hash should be equal")
}

func TestSha3_256(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Sha3_256(value)
	actual := hasher.Sha3_256(value)
	assert.Equal(t, expected, actual, "SHA3-256 hash should be equal")
}

func TestSha3_384(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Sha3_384(value)
	actual := hasher.Sha3_384(value)
	assert.Equal(t, expected, actual, "SHA3-384 hash should be equal")
}

func TestSha3_512(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Sha3_512(value)
	actual := hasher.Sha3_512(value)
	assert.Equal(t, expected, actual, "SHA3-512 hash should be equal")
}

func TestBlake2b_256(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Blake2b_256(value)
	actual := hasher.Blake2b_256(value)
	assert.Equal(t, expected, actual, "Blake2b-256 hash should be equal")
}

func TestBlake2b_512(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Blake2b_512(value)
	actual := hasher.Blake2b_512(value)
	assert.Equal(t, expected, actual, "Blake2b-512 hash should be equal")
}

func TestRipemd160(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Ripemd160(value)
	actual := hasher.Ripemd160(value)
	assert.Equal(t, expected, actual, "RIPEMD-160 hash should be equal")
}

func TestArgon2(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Argon2(value)
	actual := hasher.Argon2(value)
	assert.Equal(t, expected, actual, "Argon2 hash should be equal")
}

func TestBlake3(t *testing.T) {
	value := []byte("test data")
	expected := hasher.Blake3(value)
	actual := hasher.Blake3(value)
	assert.Equal(t, expected, actual, "Blake3 hash should be equal")
}
