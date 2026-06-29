package encrypt

import (
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
)

type StargoQueryParser struct {
	key []byte
}

func NewStargoQueryParser(key []byte, opt ...MarshalOption) *StargoQueryParser {

	return &StargoQueryParser{
		key: key,
	}
}

func (e *StargoQueryParser) Parse(target proto.Message, values url.Values, filter *utilities.DoubleArray) error {

	if v := values.Get("data"); v != "" && len(values) == 1 {

		data, err := Decode(v, string(e.key))
		if err == nil {
			if u, err := url.ParseQuery(string(data)); err == nil {
				values = u
			}
		}
	}
	return (&runtime.DefaultQueryParser{}).Parse(target, values, filter)
}
