//go:build windows && amd64

package netmgmt

import (
	"fmt"
	"os"
	"strings"
	"unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	netmgmtInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/netmgmt"
)

func AddUserToLocalGroup(user *netmgmtModelsPb.User_1, localGroup *netmgmtModelsPb.LocalGroup_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	if localGroup == nil {
		return status.Errorf(codes.InvalidArgument, "LocalGroup cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	lgrmi3_domainandname, err := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode domain and name to UTF16")
	}
	groupname, err := encode.UTF16PtrFromString(localGroup.Groupname)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode groupname to UTF16")
	}

	var bufData = &netmgmtInternalApi.LOCALGROUP_MEMBERS_INFO_3{
		Lgrmi3_domainandname: lgrmi3_domainandname,
	}
	var bufPtr = unsafe.Pointer(bufData)
	var buf = (*byte)(bufPtr)

	ret, _, _ := netmgmtInternalApi.NetLocalGroupAddMembers(
		nil, // local
		groupname,
		3, // level 3, LOCALGROUP_MEMBERS_INFO_3
		buf,
		1, // totalentries 1
	)

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_GroupNotFound:
			return status.Errorf(codes.NotFound, "Group not found")
		case win.ERROR_NO_SUCH_ALIAS:
			return status.Errorf(codes.NotFound, "Group not found")
		case win.ERROR_NO_SUCH_MEMBER:
			return status.Errorf(codes.NotFound, "User not found")
		case win.ERROR_MEMBER_IN_ALIAS:
			return status.Errorf(codes.FailedPrecondition, "User is already a member of the specified local group")
		case win.ERROR_INVALID_MEMBER:
			return status.Errorf(codes.FailedPrecondition, "User account type is invalid")
		default:
			return status.Errorf(codes.Unknown, "Failed to add user to local group (error: 0x%x)", ret)
		}
	}

	return nil
}

func RemoveUserFromLocalGroup(user *netmgmtModelsPb.User_1, localGroup *netmgmtModelsPb.LocalGroup_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	if localGroup == nil {
		return status.Errorf(codes.InvalidArgument, "LocalGroup cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	lgrmi3_domainandname, err := encode.UTF16PtrFromString(fmt.Sprintf("%s\\%s", hostname, user.Username))
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode domain and name to UTF16")
	}
	groupname, err := encode.UTF16PtrFromString(localGroup.Groupname)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode groupname to UTF16")
	}

	var bufData = &netmgmtInternalApi.LOCALGROUP_MEMBERS_INFO_3{
		Lgrmi3_domainandname: lgrmi3_domainandname,
	}
	var bufPtr = unsafe.Pointer(bufData)
	var buf = (*byte)(bufPtr)

	ret, _, _ := netmgmtInternalApi.NetLocalGroupDelMembers(
		nil, // local
		groupname,
		3, // level 3, LOCALGROUP_MEMBERS_INFO_3
		buf,
		1, // totalentries 1
	)

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_GroupNotFound:
			return status.Errorf(codes.NotFound, "Group not found")
		case win.ERROR_NO_SUCH_ALIAS:
			return status.Errorf(codes.NotFound, "Group not found")
		case win.ERROR_NO_SUCH_MEMBER:
			return status.Errorf(codes.NotFound, "User not found")
		case win.ERROR_MEMBER_NOT_IN_ALIAS:
			return status.Errorf(codes.FailedPrecondition, "User is not a member of the specified local group")
		case win.ERROR_INVALID_MEMBER:
			return status.Errorf(codes.FailedPrecondition, "User account type is invalid")
		default:
			return status.Errorf(codes.Unknown, "Failed to remove user from local group (error: 0x%x)", ret)
		}
	}

	return nil
}

func GetLocalGroups() ([]*netmgmtModelsPb.LocalGroup, error) {
	var localGroups []*netmgmtModelsPb.LocalGroup

	var buf = new(byte)
	var entriesread uint32
	var totalentries uint32
	var resumehandle uint32

	var bufDataSample netmgmtInternalApi.LOCALGROUP_INFO_0

	var ret uintptr

	for {
		ret, _, _ = netmgmtInternalApi.NetLocalGroupEnum(
			nil, // local
			0,   // level 0, LOCALGROUP_INFO_0
			&buf,
			netmgmtInternalApi.MAX_PREFERRED_LENGTH,
			&entriesread,
			&totalentries,
			&resumehandle,
		)
		if ret == netmgmtInternalApi.NERR_Success || ret == win.ERROR_MORE_DATA {
			var bufPtr = unsafe.Pointer(buf)
			for i := uint32(0); i < entriesread; i++ {
				var bufData = (*netmgmtInternalApi.LOCALGROUP_INFO_0)(unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample)))

				var localGroup = &netmgmtModelsPb.LocalGroup{
					Groupname: encode.UTF16PtrToString(bufData.Lgrpi0_name),
				}
				localGroups = append(localGroups, localGroup)
			}

			if ret == win.ERROR_MORE_DATA {
				netmgmtInternalApi.NetApiBufferFree(buf)
				buf = new(byte)
				continue
			} else {
				break
			}
		} else {
			break
		}
	}

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		default:
			return nil, status.Errorf(codes.Unknown, "Failed to get local groups (error: 0x%x)", ret)
		}
	}

	netmgmtInternalApi.NetApiBufferFree(buf)
	buf = nil

	return localGroups, nil
}

func GetUsersInLocalGroup(localGroup *netmgmtModelsPb.LocalGroup_1) ([]*netmgmtModelsPb.User_1, error) {
	var users []*netmgmtModelsPb.User_1

	if localGroup == nil {
		return nil, status.Errorf(codes.InvalidArgument, "LocalGroup cannot be nil")
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Unable to get hostname (%s)", err.Error())
	}

	groupname, err := encode.UTF16PtrFromString(localGroup.Groupname)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to encode groupname to UTF16")
	}

	var buf = new(byte)
	var entriesread uint32
	var totalentries uint32
	var resumehandle uint32

	var bufDataSample netmgmtInternalApi.LOCALGROUP_MEMBERS_INFO_3

	var ret uintptr

	for {
		ret, _, _ = netmgmtInternalApi.NetLocalGroupGetMembers(
			nil, // local
			groupname,
			3, // level 3, LOCALGROUP_MEMBERS_INFO_3
			&buf,
			netmgmtInternalApi.MAX_PREFERRED_LENGTH,
			&entriesread,
			&totalentries,
			&resumehandle,
		)
		if ret == netmgmtInternalApi.NERR_Success || ret == win.ERROR_MORE_DATA {
			var bufPtr = unsafe.Pointer(buf)
			for i := uint32(0); i < entriesread; i++ {
				var bufData = (*netmgmtInternalApi.LOCALGROUP_MEMBERS_INFO_3)(unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample)))

				domainandname := encode.UTF16PtrToString(bufData.Lgrmi3_domainandname)
				domainandname_splitted := strings.Split(domainandname, "\\")
				if len(domainandname_splitted) != 2 {
					continue
				}

				domain := domainandname_splitted[0]
				if domain != hostname {
					continue
				}

				name := domainandname_splitted[1]

				var user = &netmgmtModelsPb.User_1{
					Username: name,
				}
				users = append(users, user)
			}

			if ret == win.ERROR_MORE_DATA {
				netmgmtInternalApi.NetApiBufferFree(buf)
				buf = new(byte)
				continue
			} else {
				break
			}
		} else {
			break
		}
	}

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_GroupNotFound:
			return nil, status.Errorf(codes.NotFound, "Group not found")
		default:
			return nil, status.Errorf(codes.Unknown, "Failed to get users in local group (error: 0x%x)", ret)
		}
	}

	netmgmtInternalApi.NetApiBufferFree(buf)
	buf = nil

	return users, nil
}
