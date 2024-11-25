package api

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func MetaString(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	value, ok := md[key]
	if !ok {
		return ""
	}
	return value[0]
}

func MetaInt(ctx context.Context, key string) int {
	v := MetaString(ctx, key)
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return r
}

func MetaHost(ctx context.Context) string {
	return MetaString(ctx, "x-forwarded-host")
}

func MetaIp(ctx context.Context) string {
	return MetaString(ctx, "x-forwarded-for")
}
func MetaMethod(ctx context.Context) string {
	return MetaString(ctx, "g-method")
}
func MetaToken(ctx context.Context) string {
	return MetaString(ctx, "token")
}
