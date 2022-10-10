//go:build windows && amd64

package main

import (
	"context"
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/beevik/ntp"
	"github.com/go-co-op/gocron"
	"github.com/golang-jwt/jwt/v4"

	fileioServicePb "github.com/s77rt/rdpcloud/proto/go/services/fileio"
	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
	shellServicePb "github.com/s77rt/rdpcloud/proto/go/services/shell"
	sysinfoServicePb "github.com/s77rt/rdpcloud/proto/go/services/sysinfo"
	termservServicePb "github.com/s77rt/rdpcloud/proto/go/services/termserv"
	"github.com/s77rt/rdpcloud/server/go/config"
	customJWT "github.com/s77rt/rdpcloud/server/go/jwt"
	fileioService "github.com/s77rt/rdpcloud/server/go/services/fileio"
	netmgmtService "github.com/s77rt/rdpcloud/server/go/services/netmgmt"
	secauthnService "github.com/s77rt/rdpcloud/server/go/services/secauthn"
	secauthzService "github.com/s77rt/rdpcloud/server/go/services/secauthz"
	shellService "github.com/s77rt/rdpcloud/server/go/services/shell"
	sysinfoService "github.com/s77rt/rdpcloud/server/go/services/sysinfo"
	termservService "github.com/s77rt/rdpcloud/server/go/services/termserv"
)

var (
	Version string = "dev"

	ServerName string
	ServerIP   string

	LicenseExpDate string // YYYY-MM-DD format
)

func init() {
	// Check Server IP
	conn, err := net.Dial("udp4", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Failed to dial udp4: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	if !(localAddr.IP.Equal(net.ParseIP(ServerIP))) {
		log.Fatalf("Server IP does not match; expected %s, got %s", ServerIP, localAddr.IP)
	}

	// Check License Exp. Date
	if LicenseExpDate != "" {
		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Day().Do(func() {
			timeNow, err := ntp.Time("time.google.com")
			if err != nil {
				log.Fatalf("Failed to read time: %v", err)
			}

			timeExp, err := time.Parse("2006-01-02", LicenseExpDate)
			if err != nil {
				log.Fatalf("Failed to read exp time: %v", err)
			}

			if timeNow.After(timeExp) {
				log.Fatalf("License expired on %s", timeExp)
			}
		})
		s.StartAsync()
	}
}

//go:embed cert
var c embed.FS

func main() {
	var port int
	flag.IntVar(&port, "port", 5027, "port on which the server will listen")
	flag.Parse()

	log.Printf("Running RDPCloud Server (Version: %s)", Version)

	if LicenseExpDate == "" {
		log.Printf("Licensed to %s (%s)", ServerName, ServerIP)
	} else {
		log.Printf("Licensed to %s (%s) [Exp. Date: %s]", ServerName, ServerIP, LicenseExpDate)
	}

	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %v", lis.Addr())

	serverCert, err := c.ReadFile("cert/server-cert.pem")
	if err != nil {
		log.Fatalf("Failed to read server-cert pem file: %v", err)
	}
	serverKey, err := c.ReadFile("cert/server-key.pem")
	if err != nil {
		log.Fatalf("Failed to read server-key pem file: %v", err)
	}

	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		log.Fatalf("Failed to load tls key pair: %v", err)
	}

	grpcOpts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		grpc.UnaryInterceptor(ensureValidToken), // Currently we only use unary RPCs

		// Enable TLS for all incoming connections
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}

	s := grpc.NewServer(grpcOpts...)
	netmgmtServicePb.RegisterNetmgmtServer(s, &netmgmtService.Server{})
	secauthnServicePb.RegisterSecauthnServer(s, &secauthnService.Server{})
	secauthzServicePb.RegisterSecauthzServer(s, &secauthzService.Server{})
	fileioServicePb.RegisterFileioServer(s, &fileioService.Server{})
	termservServicePb.RegisterTermservServer(s, &termservService.Server{})
	shellServicePb.RegisterShellServer(s, &shellService.Server{})
	sysinfoServicePb.RegisterSysinfoServer(s, &sysinfoService.Server{})

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
