package testutil

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func NewGrpcServer(t *testing.T, setupFn func(*grpc.Server)) (*grpc.Server, *grpc.ClientConn) {
	// create a in-memory buffer listener
	lis := bufconn.Listen(1024 * 1024)

	// create the gRPC server
	s := grpc.NewServer()

	// run the setup function
	if setupFn != nil {
		setupFn(s)
	}

	// start the gRPC server in the background
	go func() {
		if err := s.Serve(lis); err != nil {
			// TODO: can we somehow use t.Fatal here?
			// Problem: t.Fatal has to be run in the same goroutine as the test itself.
			log.Fatalf("server exited with error: %v", err)
		}
	}()

	// create a dialer that connects to the in-memory buffer listener
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	// create a gRPC client that connects to the mock server
	c, err := grpc.NewClient(
		"localhost",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return s, c
}
