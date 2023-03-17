package server

import (
	"context"
	pb "user/pkg/pb"
)

func (e *handler) FetchUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return &pb.User{
		Uid:  123,
		Name: "hello stargo",
	}, nil
}
