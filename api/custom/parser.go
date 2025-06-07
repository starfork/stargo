package custom

import (
	"encoding/base64"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"google.golang.org/protobuf/proto"
)

type StargoQueryParser struct {
}

// Parse parses query parameters and populates the appropriate fields in the gRPC request message.
func (p *StargoQueryParser) Parse(target proto.Message, values url.Values, filter *utilities.DoubleArray) error {

	if v := values.Get("data"); v != "" && len(values) == 1 {
		data, err := base64.StdEncoding.DecodeString(v)
		if err == nil {
			if u, err := url.ParseQuery(string(data)); err == nil {
				values = u
			}
		}
	}
	return (&runtime.DefaultQueryParser{}).Parse(target, values, filter)
}
