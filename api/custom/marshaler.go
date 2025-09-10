package custom

import (
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/starfork/stargo/api"
	"google.golang.org/protobuf/proto"
)

type StargoMarshaler struct {
	runtime.JSONPb
	key  []byte
	opts *MarshalOptions
}

func NewStargoMarshaler(key []byte, opt ...MarshalOption) *StargoMarshaler {

	m := &StargoMarshaler{
		key: key,
	}
	for _, o := range opt {
		o(m.opts)
	}

	m.JSONPb = api.DefaultMarshalerOption
	return m
}

func (e *StargoMarshaler) Marshal(v any) ([]byte, error) {
	fmt.Println(v)
	jsonData, err := e.JSONPb.Marshal(v.(proto.Message))
	if err != nil {
		return nil, err
	}
	//fmt.Println("marshaler")
	//encoded := base64.StdEncoding.EncodeToString(append(e.prefix, jsonData...))

	return Encode(e.key, jsonData)
}

func (e *StargoMarshaler) ContentType(_ any) string {
	return "text/plain"
}
