package custom

import (
	"encoding/base64"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

type StargoMarshaler struct {
	runtime.JSONPb
}

func (m *StargoMarshaler) Marshal(v any) ([]byte, error) {
	jsonData, err := m.JSONPb.Marshal(v.(proto.Message))
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(jsonData)
	return []byte(encoded), nil
}

func (m *StargoMarshaler) ContentType(_ any) string {
	return "text/plain"
}
