package request

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

func Meta(ctx context.Context, r *http.Request) metadata.MD {
	md := make(map[string]string)
	//fmt.Println(r.Header)
	maps := map[string]string{
		"access-token": "Access-Token",
		"device-id":    "Device-Id",
		"device":       "Device",
		"version":      "Version",
	}
	for k, v := range maps {
		if value, ok := r.Header[v]; ok {
			md[k] = value[0]
		}
	}
	if len(md) > 0 {
		return metadata.New(md)
	}
	return nil
}
