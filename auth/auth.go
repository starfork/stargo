package auth

import (
	"context"
)

type Auth interface {
	Validate(token string) error
}

// get token from ctx
func GetToken(ctx context.Context) string {

	return ""
}
