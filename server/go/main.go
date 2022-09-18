//go:build windows && amd64

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/golang-jwt/jwt/v4"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
	"github.com/s77rt/rdpcloud/server/go/config"
	customJWT "github.com/s77rt/rdpcloud/server/go/jwt"
	netmgmtService "github.com/s77rt/rdpcloud/server/go/services/netmgmt"
	secauthnService "github.com/s77rt/rdpcloud/server/go/services/secauthn"
	secauthzService "github.com/s77rt/rdpcloud/server/go/services/secauthz"
)

var (
	Version string = "dev"
)

func main() {
	log.Printf("Running RDPCloud Server (Version: %s)", Version)

	lis, err := net.Listen("tcp", ":5027")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %v", lis.Addr())

	grpcOpts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken), // Currently we only use unary RPCs

		// Enable TLS for all incoming connections
		// grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	s := grpc.NewServer(grpcOpts...)
	netmgmtServicePb.RegisterNetmgmtServer(s, &netmgmtService.Server{})
	secauthnServicePb.RegisterSecauthnServer(s, &secauthnService.Server{})
	secauthzServicePb.RegisterSecauthzServer(s, &secauthzService.Server{})

	// Register reflection service
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	accessLevel, found := config.AceessLevel[info.FullMethod]
	if !found {
		return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to execute the specified operation")
	}

	if accessLevel == 0 {
		// Skip authentication
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "Missing metadata")
	}

	authorization := md["authorization"]
	if len(authorization) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	tokenString := strings.TrimPrefix(authorization[0], "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &customJWT.UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect

		signingMethod, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if signingMethod != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.Secret, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token (%s)", err.Error())
	}

	if claims, ok := token.Claims.(*customJWT.UserClaims); ok && token.Valid {
		if claims.Privilege < accessLevel {
			return nil, status.Errorf(codes.PermissionDenied, "You don't have permission to execute the specified operation")
		}
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	// Continue execution of handler after ensuring a valid token
	return handler(ctx, req)
}
