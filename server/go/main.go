//go:build windows && amd64

package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
	netmgmtService "github.com/s77rt/rdpcloud/server/go/services/netmgmt"
)

func main() {
	lis, err := net.Listen("tcp", ":5027")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	netmgmtServicePb.RegisterNetmgmtServer(s, &netmgmtService.Server{})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
