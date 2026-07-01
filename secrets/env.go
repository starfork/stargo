package secrets

import (
	"os"
)

func init() {
	Register("env", func(c *Config) (SecretManager, error) {
		return &EnvManager{}, nil
	})
}

type EnvManager struct{}

func (e *EnvManager) Resolve(reference string) (string, error) {
	ref, err := ParseReference(reference)
	if err != nil {
		return "", err
	}
	return os.Getenv(ref.Path), nil
}

func (e *EnvManager) Close() error { return nil }
