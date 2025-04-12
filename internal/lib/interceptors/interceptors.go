package interceptors

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
)

func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)
		return handler(ctx, req)
	}
}

func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
