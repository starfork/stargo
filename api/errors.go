package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type APIError struct {
	Service string        `json:"service"`
	Method  string        `json:"method"`
	Code    string        `json:"code"`
	Msg     string        `json:"msg"`
	Details []interface{} `json:"details,omitempty"`
}

func StargoHTTPError(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	apiError := APIError{
		Service: "UnknownService",
		Method:  r.Method,
		Code:    st.Code().String(),
		Msg:     st.Message(),
	}

	for _, d := range st.Details() {
		switch info := d.(type) {
		case *errdetails.ErrorInfo:
			apiError.Service = info.Metadata["Service"]
			apiError.Method = info.Metadata["Method"]
		default:
			apiError.Details = append(apiError.Details, info)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))
	json.NewEncoder(w).Encode(apiError)
}
