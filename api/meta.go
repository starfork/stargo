package api

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

func MetaString(ctx context.Context, key string, out ...bool) string {
	var md metadata.MD
	var ok bool
	if len(out) > 0 && out[0] {
		md, ok = metadata.FromOutgoingContext(ctx)
	} else {
		md, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return ""
	}
	key = strings.ToLower(key)
	value, ok := md[key]
	if !ok {
		return ""
	}
	return value[0]
}

func MetaInt(ctx context.Context, key string, out ...bool) int {
	v := MetaString(ctx, key, out...)
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return r
}

func MetaHost(ctx context.Context, out ...bool) string {
	return MetaString(ctx, META_HOST, out...)
}

func MetaIp(ctx context.Context, out ...bool) string {
	return MetaString(ctx, META_IP, out...)
}
func MetaFp(ctx context.Context, out ...bool) string {
	return MetaString(ctx, META_FP, out...)
}
func MetaMethod(ctx context.Context, out ...bool) string {
	return MetaString(ctx, META_METHOD, out...)
}
func MetaToken(ctx context.Context, out ...bool) string {
	return MetaString(ctx, META_TOKEN, out...)
}

func MetaLang(ctx context.Context, out ...bool) string {
	//zh-CN,zh;q=0.9,en;q=0.8
	str := MetaString(ctx, META_LANG, out...)
	if str != "" {
		tmp := strings.Split(str, ",")
		return tmp[0]
	}
	return ""

}
