package crypt

type Registry interface {
	Init(...Option) error
	Encrypt() ([]byte, error)
	Decrypt() ([]byte, error)
}

type Option func(*Options)
