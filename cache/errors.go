package cache

import "errors"

var (
	ErrKeyExpired = errors.New("key expired")

	ErrMultiGetFailed = errors.New("multi get failed")

	ErrIncrementOverflow = errors.New("incream overflow")

	ErrDecrementOverflow = errors.New("cecrement overflow")

	ErrNotIntegerType = errors.New("not integer type")
)
