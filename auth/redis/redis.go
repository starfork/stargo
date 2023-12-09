package jwt

import (
	"github.com/starfork/stargo/auth"
)

type RedisAuth struct{}

func New() auth.Auth {
	return &RedisAuth{}
}

func (e *RedisAuth) Validate(token string) error {
	return nil
}
