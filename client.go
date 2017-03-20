package main

import (
	"context"

	"google.golang.org/grpc"
)

type client struct {
	conn      *grpc.ClientConn
	stream    Streamer_StreamClient
	cancelCtx context.CancelFunc
}

func newClient(addr string) (*client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	sc := NewStreamerClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	stream, err := sc.Stream(ctx, &StreamRequest{})
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &client{conn, stream, cancel}, nil
}

// Close asynchronously shuts down the connection.
func (c *client) Close() error {
	c.cancelCtx()
	c.conn.Close()
	return nil
}

func (c *client) Recv() (string, error) {
	rsp, err := c.stream.Recv()
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}
