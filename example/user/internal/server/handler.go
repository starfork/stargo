package server

import (
	pb "user/pkg/pb"
)

type handler struct {
	pb.UnimplementedUserHandlerServer
}

func New() *handler {
	return &handler{}
}
