package types

import "errors"

var (
	ErrKeyNotFound  = errors.New("Key not found")
	ErrTypeMismatch = errors.New("Type mismatch")
)
