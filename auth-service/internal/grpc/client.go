package grpcclient

import (
	"context"
	"fmt"

	pb "amelli/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.AlhelisServiceClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create user-service grpc client: %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewAlhelisServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) CreateProfile(ctx context.Context, userID int64) error {
	_, err := c.client.CreateProfile(
		ctx,
		&pb.CreateProfileRequest{
			UserId: fmt.Sprintf("%d", userID),
		},
	)
	if err != nil {
		return fmt.Errorf("create profile via user-service: %w", err)
	}

	return nil
}
