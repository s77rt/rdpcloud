package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
)

func main() {
	conn, err := grpc.Dial("51.89.161.169:5027", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	netmgmtClient := netmgmtServicePb.NewNetmgmtClient(conn)

	getUsersResponse, err := netmgmtClient.GetUsers(context.Background(), &netmgmtServicePb.GetUsersRequest{})
	if err != nil {
		fmt.Println("error", err)
		return
	}
	for _, user := range getUsersResponse.Users {
		fmt.Println("Username:", user.Username)
		fmt.Println("Privilege:", user.Privilege)
	}

	user := &netmgmtModelsPb.User{
		Username: "test",
		Password: "testyocoolawesome",
	}
	_, err = netmgmtClient.AddUser(context.Background(), &netmgmtServicePb.AddUserRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error", err)
		return
	}
	fmt.Println("User Added!")
}
