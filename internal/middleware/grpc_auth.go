package middleware

import (
	"context"
	"log"
	"strings"

	"github.com/alexe0110/chat-system/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if info.FullMethod == "/chat.UserService/GetUser" {
			return handler(ctx, req)
		}

		reqMetadata, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata isn't get")
		}
		token := reqMetadata.Get("authorization")

		if len(token) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "missing token")
		}

		onlyToken := strings.TrimPrefix(token[0], "Bearer ")

		userID, err := auth.ParseToken(onlyToken, secret)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token is not verified")
		}

		log.Println("Success login user", userID)
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func AuthStreamInterceptor(secret string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if info.FullMethod == "/chat.UserService/GetUser" {
			return handler(srv, ss)
		}

		ctx := ss.Context()
		reqMetadata, ok := metadata.FromIncomingContext(ctx)

		if !ok {
			return status.Errorf(codes.Unauthenticated, "metadata isn't get")
		}
		token := reqMetadata.Get("authorization")

		if len(token) == 0 {
			return status.Errorf(codes.Unauthenticated, "missing token")
		}

		onlyToken := strings.TrimPrefix(token[0], "Bearer ")

		userID, err := auth.ParseToken(onlyToken, secret)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "token is not verified")
		}

		log.Println("Success login user", userID)
		err = handler(srv, ss)
		return err
	}

}
