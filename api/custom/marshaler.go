package custom

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

type StargoMarshaler struct {
	runtime.JSONPb

	key []byte
}

func NewStargoMarshaler(key ...[]byte) *StargoMarshaler {
	k := Key
	if len(key) > 0 {
		k = key[0]
	}
	return &StargoMarshaler{
		key: k,
	}
}

func (e *StargoMarshaler) Marshal(v any) ([]byte, error) {
	jsonData, err := e.JSONPb.Marshal(v.(proto.Message))
	if err != nil {
		return nil, err
	}
	//encoded := base64.StdEncoding.EncodeToString(append(e.prefix, jsonData...))

	return Encode(e.key, jsonData)
}

func (e *StargoMarshaler) ContentType(_ any) string {
	return "text/plain"
}
