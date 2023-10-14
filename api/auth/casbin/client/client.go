package client

import (
	"context"

	pb "github.com/casbin/casbin-server/proto"
	"google.golang.org/grpc"
)

// Client is a wrapper around proto.CasbinClient, and can be used to create an Enforcer.
type Client struct {
	remoteClient pb.CasbinClient
}

// NewClient creates and returns a new client for casbin-server.
func NewClient(ctx context.Context, address string, opts ...grpc.DialOption) (*Client, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, err
	}
	c := pb.NewCasbinClient(conn)

	return &Client{
		remoteClient: c,
	}, nil
}
