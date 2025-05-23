package auth

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	headerAuthorize = "authorization"
)

// AuthFromMD is a helper function for extracting the :authorization header from the gRPC metadata of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with gRPC status `Unauthenticated` is returned.
func AuthFromMD(ctx context.Context, expectedScheme string) (string, error) {
	val := metadata.ExtractIncoming(ctx).Get(headerAuthorize)
	if val == "" {
		return "", status.Errorf(codes.Unauthenticated, "%s", "request unauthenticated with "+expectedScheme)
	}
	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", status.Errorf(codes.Unauthenticated, "bad authorization string")
	}
	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", status.Errorf(codes.Unauthenticated, "%s", "request unauthenticated with "+expectedScheme)
	}
	return splits[1], nil
}
