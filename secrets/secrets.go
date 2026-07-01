package secrets

import (
	"fmt"
	"strings"
)

type Reference struct {
	Prefix string
	Path   string
	Key    string
}

func ParseReference(s string) (*Reference, error) {
	idx := strings.Index(s, ":")
	if idx < 0 {
		return nil, fmt.Errorf("secrets: invalid reference syntax: %s", s)
	}
	prefix := s[:idx]
	rest := s[idx+1:]

	hashIdx := strings.LastIndex(rest, "#")
	if hashIdx < 0 {
		return &Reference{Prefix: prefix, Path: rest, Key: ""}, nil
	}

	return &Reference{
		Prefix: prefix,
		Path:   rest[:hashIdx],
		Key:    rest[hashIdx+1:],
	}, nil
}

func IsReference(s string) bool {
	_, err := ParseReference(s)
	return err == nil
}

type SecretManager interface {
	Resolve(reference string) (string, error)
	Close() error
}

var factories = make(map[string]func(*Config) (SecretManager, error))

func Register(name string, factory func(*Config) (SecretManager, error)) {
	factories[name] = factory
}

func New(name string, conf *Config) (SecretManager, error) {
	if f, ok := factories[name]; ok {
		return f(conf)
	}
	return nil, fmt.Errorf("secrets: unknown driver %q", name)
}
