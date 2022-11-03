//go:build windows && amd64

package main

import (
	"crypto/tls"
	"embed"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"github.com/beevik/ntp"
	externalip "github.com/glendc/go-external-ip"
	"github.com/go-co-op/gocron"

	fileioServicePb "github.com/s77rt/rdpcloud/proto/go/services/fileio"
	msiServicePb "github.com/s77rt/rdpcloud/proto/go/services/msi"
	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
	shellServicePb "github.com/s77rt/rdpcloud/proto/go/services/shell"
	shutdownServicePb "github.com/s77rt/rdpcloud/proto/go/services/shutdown"
	sysinfoServicePb "github.com/s77rt/rdpcloud/proto/go/services/sysinfo"
	termservServicePb "github.com/s77rt/rdpcloud/proto/go/services/termserv"
	"github.com/s77rt/rdpcloud/server/go/auth"
	"github.com/s77rt/rdpcloud/server/go/license"
	fileioService "github.com/s77rt/rdpcloud/server/go/services/fileio"
	msiService "github.com/s77rt/rdpcloud/server/go/services/msi"
	netmgmtService "github.com/s77rt/rdpcloud/server/go/services/netmgmt"
	secauthnService "github.com/s77rt/rdpcloud/server/go/services/secauthn"
	secauthzService "github.com/s77rt/rdpcloud/server/go/services/secauthz"
	shellService "github.com/s77rt/rdpcloud/server/go/services/shell"
	shutdownService "github.com/s77rt/rdpcloud/server/go/services/shutdown"
	sysinfoService "github.com/s77rt/rdpcloud/server/go/services/sysinfo"
	termservService "github.com/s77rt/rdpcloud/server/go/services/termserv"
)

var (
	Version string = "dev"

	licenseInfo *license.License
)

func init() {
	// Read License
	var err error
	licenseInfo, err = license.Read()
	if err != nil {
		log.Fatalf("Failed to read license: %v", err)
	}

	// Check License Server Local IP
	conn, err := net.Dial("udp4", "8.8.8.8:53")
	if err != nil {
		log.Fatalf("Failed to dial udp4: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	if !(localAddr.IP.Equal(licenseInfo.ServerLocalIP)) {
		log.Fatalf("Server Local IP does not match; expected %s, got %s", licenseInfo.ServerLocalIP, localAddr.IP)
	}

	// Check License Server Public IP
	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(4)

	publicIP, err := consensus.ExternalIP()
	if err != nil {
		log.Fatalf("Failed to read public IP: %v", err)
	}

	if !(publicIP.Equal(licenseInfo.ServerPublicIP)) {
		log.Fatalf("Server Public IP does not match; expected %s, got %s", licenseInfo.ServerPublicIP, publicIP)
	}

	// Check License Exp. Date
	if !licenseInfo.ExpDate.IsZero() {
		s := gocron.NewScheduler(time.UTC)
		s.Every(1).Day().Do(func() {
			timeNow, err := ntp.Time("time.google.com")
			if err != nil {
				log.Fatalf("Failed to read time: %v", err)
			}

			if timeNow.After(licenseInfo.ExpDate) {
				log.Fatalf("License expired on %s", licenseInfo.ExpDate)
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

	if licenseInfo.ExpDate.IsZero() {
		log.Printf("Licensed to %s (%s | %s)", licenseInfo.ServerName, licenseInfo.ServerLocalIP, licenseInfo.ServerPublicIP)
	} else {
		log.Printf("Licensed to %s (%s | %s) [Exp. Date: %s]", licenseInfo.ServerName, licenseInfo.ServerLocalIP, licenseInfo.ServerPublicIP, licenseInfo.ExpDate)
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
		// Unary interceptors
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			auth.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		)),

		// Stream interceptors
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			auth.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		)),

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
	shutdownServicePb.RegisterShutdownServer(s, &shutdownService.Server{})
	msiServicePb.RegisterMsiServer(s, &msiService.Server{})

	// Register reflection service
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
