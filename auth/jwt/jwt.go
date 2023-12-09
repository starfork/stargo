package jwt

import (
	"github.com/starfork/stargo/auth"
)

type JwtAuth struct{}

func New() auth.Auth {
	return &JwtAuth{}
}

func (e *JwtAuth) Validate(token string) error {
	return nil
}
