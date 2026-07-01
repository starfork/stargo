package secrets

func init() {
	Register("plain", func(c *Config) (SecretManager, error) {
		return &PlainManager{}, nil
	})
}

type PlainManager struct{}

func (p *PlainManager) Resolve(reference string) (string, error) {
	ref, err := ParseReference(reference)
	if err != nil {
		return "", err
	}
	return ref.Path, nil
}

func (p *PlainManager) Close() error { return nil }
