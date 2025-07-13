package grpcutil

import (
	"google.golang.org/grpc/codes"
)

func ShouldRetryCode(code codes.Code) bool {
	return code == codes.Unavailable || code == codes.DeadlineExceeded
}
