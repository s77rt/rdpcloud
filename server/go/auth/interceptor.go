//go:build windows && amd64

package auth

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/s77rt/rdpcloud/server/go/config"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		accessLevel, found := config.AccessLevel[info.FullMethod]
		if !found {
			return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to execute the specified operation")
		}

		if accessLevel == 0 {
			// Skip authentication
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
		}

		authorization := md["authorization"]
		if len(authorization) < 1 {
			return nil, status.Errorf(codes.Unauthenticated, "Missing token")
		}

		tokenString := strings.TrimPrefix(authorization[0], "Bearer ")

		token, err := ParseTokenWithUserClaims(tokenString)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Unable to parse token (%s)", err.Error())
		}

		if !token.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		}

		userClaims, err := TokenUserClaims(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token user claims (%s)", err.Error())
		}

		if userClaims.Privilege < accessLevel {
			return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to execute the specified operation")
		}

		newCtx := NewContextWithUserClaims(ctx, userClaims)

		// Continue execution of handler after ensuring a valid token
		return handler(newCtx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Allow server reflection
		if info.FullMethod == "/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo" {
			return handler(srv, ss)
		}

		return status.Errorf(codes.Unimplemented, "Stream RPCs are not supported")
	}
}
