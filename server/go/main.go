//go:build windows && amd64

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
	netmgmtService "github.com/s77rt/rdpcloud/server/go/services/netmgmt"
	secauthnService "github.com/s77rt/rdpcloud/server/go/services/secauthn"
)

func main() {
	lis, err := net.Listen("tcp", ":5027")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %v", lis.Addr())

	s := grpc.NewServer()
	netmgmtServicePb.RegisterNetmgmtServer(s, &netmgmtService.Server{})
	secauthnServicePb.RegisterSecauthnServer(s, &secauthnService.Server{})

	// Register reflection service
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
