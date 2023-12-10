package marshaler

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pj "google.golang.org/protobuf/encoding/protojson"
)

var MarshalerOptions = runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
	Marshaler: &runtime.JSONPb{
		MarshalOptions: pj.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: pj.UnmarshalOptions{
			DiscardUnknown: true,
		},
	},
})
