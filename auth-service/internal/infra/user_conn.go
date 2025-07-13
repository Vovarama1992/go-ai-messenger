package grpc

import (
	"context"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

// NewUserServiceConn возвращает gRPC соединение с retry и timeout логикой.
func NewUserServiceConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	return grpc.DialContext(
		ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(
				grpc_retry.WithMax(3),
				grpc_retry.WithPerRetryTimeout(500*time.Millisecond),
				grpc_retry.WithCodes(codes.Unavailable, codes.DeadlineExceeded),
			),
		),
	)
}
