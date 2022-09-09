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

var netmgmtClient netmgmtServicePb.NetmgmtClient

func AddUser() {
	fmt.Println("AddUser")
	user := &netmgmtModelsPb.User{
		Username: "test",
		Password: "cai@6Hae+vohniut",
	}
	_, err := netmgmtClient.AddUser(context.Background(), &netmgmtServicePb.AddUserRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("User Added!")
	}
}

func DeleteUser() {
	fmt.Println("DeleteUser")
	user := &netmgmtModelsPb.User{
		Username: "test",
	}
	_, err := netmgmtClient.DeleteUser(context.Background(), &netmgmtServicePb.DeleteUserRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("User Deleted!")
	}
}

func GetUsers() {
	fmt.Println("GetUsers")
	resp, err := netmgmtClient.GetUsers(context.Background(), &netmgmtServicePb.GetUsersRequest{})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, user := range resp.Users {
			fmt.Println("Username:", user.Username)
			fmt.Println("Privilege:", user.Privilege)
			fmt.Println("Flags:", user.Flags)
		}
	}
}

func GetUser() {
	fmt.Println("GetUser")
	user := &netmgmtModelsPb.User{
		Username: "test",
	}
	resp, err := netmgmtClient.GetUser(context.Background(), &netmgmtServicePb.GetUserRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("Username:", resp.User.Username)
		fmt.Println("Privilege:", resp.User.Privilege)
		fmt.Println("Flags:", resp.User.Flags)
	}
}

func GetUserLocalGroups() {
	fmt.Println("GetUserLocalGroups")
	user := &netmgmtModelsPb.User{
		Username: "test",
	}
	resp, err := netmgmtClient.GetUserLocalGroups(context.Background(), &netmgmtServicePb.GetUserLocalGroupsRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, localGroup := range resp.LocalGroups {
			fmt.Println("Groupname:", localGroup.Groupname)
		}
	}
}

func ChangeUserPassword() {
	fmt.Println("ChangeUserPassword")
	user := &netmgmtModelsPb.User{
		Username: "test",
		Password: "cai@6Hae+vohniut",
	}
	_, err := netmgmtClient.ChangeUserPassword(context.Background(), &netmgmtServicePb.ChangeUserPasswordRequest{
		User: user,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("User Password Changed!")
	}
}

func AddUserToLocalGroup(groupname string) {
	fmt.Println("AddUserToLocalGroup")
	user := &netmgmtModelsPb.User{
		Username: "test",
	}
	localGroup := &netmgmtModelsPb.LocalGroup{
		Groupname: groupname,
	}
	_, err := netmgmtClient.AddUserToLocalGroup(context.Background(), &netmgmtServicePb.AddUserToLocalGroupRequest{
		User:       user,
		LocalGroup: localGroup,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("User Added To Local Group!")
	}
}

func RemoveUserFromLocalGroup() {
	fmt.Println("RemoveUserFromLocalGroup")
	user := &netmgmtModelsPb.User{
		Username: "test",
	}
	localGroup := &netmgmtModelsPb.LocalGroup{
		Groupname: "Users",
	}
	_, err := netmgmtClient.RemoveUserFromLocalGroup(context.Background(), &netmgmtServicePb.RemoveUserFromLocalGroupRequest{
		User:       user,
		LocalGroup: localGroup,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("User Removed From Local Group!")
	}
}

func GetLocalGroups() {
	fmt.Println("GetLocalGroups")
	resp, err := netmgmtClient.GetLocalGroups(context.Background(), &netmgmtServicePb.GetLocalGroupsRequest{})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, localGroup := range resp.LocalGroups {
			fmt.Println("Groupname:", localGroup.Groupname)
		}
	}
}

func GetUsersInLocalGroup() {
	fmt.Println("GetUsersInLocalGroup")
	localGroup := &netmgmtModelsPb.LocalGroup{
		Groupname: "Users",
	}
	resp, err := netmgmtClient.GetUsersInLocalGroup(context.Background(), &netmgmtServicePb.GetUsersInLocalGroupRequest{
		LocalGroup: localGroup,
	})
	if err != nil {
		fmt.Println("error:", err)
	} else {
		for _, user := range resp.Users {
			fmt.Println("Username:", user.Username)
		}
	}
}

func main() {
	conn, err := grpc.Dial("51.89.161.169:5027", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	netmgmtClient = netmgmtServicePb.NewNetmgmtClient(conn)

	GetUsers()
	AddUser()
	AddUserToLocalGroup("Users")
	AddUserToLocalGroup("Remote Desktop Users")
	GetUser()
}
