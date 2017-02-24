package main

import (
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
)

type testStreamer struct {
	messages chan string
}

func (s *testStreamer) Stream(_ *StreamRequest, stream Streamer_StreamServer) error {
	for {
		select {
		case msg := <-s.messages:
			if err := stream.Send(&StreamResponse{Message: msg}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

func TestStreamCloseBeforeStop(t *testing.T) {
	srv, streamer, addr := startServer(t)
	defer srv.Stop()

	c, err := newClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	messages := make(chan string)
	go receiveMessages(c, messages)

	streamer.messages <- "a"
	if got := <-messages; got != "a" {
		t.Fatalf("got %v, want %v", got, "a")
	}

	c.Close()
	srv.GracefulStop()
}

func TestStreamCloseAfterStop(t *testing.T) {
	srv, streamer, addr := startServer(t)
	defer srv.Stop()

	c, err := newClient(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	messages := make(chan string)
	go receiveMessages(c, messages)

	streamer.messages <- "a"
	if got := <-messages; got != "a" {
		t.Fatalf("got %v, want %v", got, "a")
	}

	go func() {
		time.Sleep(time.Second)
		c.Close()
	}()

	srv.GracefulStop()
}

func startServer(t *testing.T) (*grpc.Server, *testStreamer, string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	srv := grpc.NewServer()

	streamer := &testStreamer{messages: make(chan string)}
	RegisterStreamerServer(srv, streamer)

	go srv.Serve(ln)

	addr := ln.Addr().String()

	return srv, streamer, addr
}

func receiveMessages(c *client, msgs chan string) {
	for {
		msg, err := c.Recv()
		if err != nil {
			return
		}

		msgs <- msg
	}
}
