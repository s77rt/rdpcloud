//go:build windows && amd64

package netmgmt

import (
	"unsafe"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	"github.com/s77rt/rdpcloud/server/go/internal/encode"
	"github.com/s77rt/rdpcloud/server/go/internal/secure"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	netmgmtInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/netmgmt"
)

func AddUser(user *netmgmtModelsPb.User_3) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	usri2_name, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}
	usri2_password, err := encode.UTF16PtrFromString(user.Password)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode password to UTF16")
	}

	var bufData = &netmgmtInternalApi.USER_INFO_2{
		Usri2_name:         usri2_name,
		Usri2_password:     usri2_password,
		Usri2_priv:         netmgmtInternalApi.USER_PRIV_USER,
		Usri2_flags:        netmgmtInternalApi.UF_NORMAL_ACCOUNT | netmgmtInternalApi.UF_SCRIPT | netmgmtInternalApi.UF_DONT_EXPIRE_PASSWD,
		Usri2_acct_expires: netmgmtInternalApi.TIMEQ_FOREVER,
	}
	var bufPtr = unsafe.Pointer(bufData)
	var buf = (*byte)(bufPtr)

	var parm_err uint32

	ret, _, _ := netmgmtInternalApi.NetUserAdd(
		nil, // local
		2,   // level 2, USER_INFO_2
		buf,
		&parm_err,
	)

	user.Password = ""
	secure.ZeroMemoryUint16FromPtr(usri2_password)
	usri2_password = nil

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_BadPassword:
			return status.Errorf(codes.InvalidArgument, "Bad password")
		case netmgmtInternalApi.NERR_UserExists:
			return status.Errorf(codes.AlreadyExists, "User already exists")
		case netmgmtInternalApi.NERR_GroupExists:
			return status.Errorf(codes.AlreadyExists, "Group already exists")
		case netmgmtInternalApi.NERR_PasswordTooShort:
			return status.Errorf(codes.InvalidArgument, "Password does not meet the password policy requirements")
		default:
			return status.Errorf(codes.Unknown, "Failed to add user (error: 0x%x)", ret)
		}
	}

	return nil
}

func DeleteUser(user *netmgmtModelsPb.User_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	username, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}

	ret, _, _ := netmgmtInternalApi.NetUserDel(
		nil, // local
		username,
	)

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_UserNotFound:
			return status.Errorf(codes.NotFound, "User not found")
		default:
			return status.Errorf(codes.Unknown, "Failed to delete user (error: 0x%x)", ret)
		}
	}

	return nil
}

func GetUsers() ([]*netmgmtModelsPb.User, error) {
	var users []*netmgmtModelsPb.User

	var buf = new(byte)
	var entriesread uint32
	var totalentries uint32
	var resumehandle uint32

	var bufDataSample netmgmtInternalApi.USER_INFO_2

	var ret uintptr

	for {
		ret, _, _ = netmgmtInternalApi.NetUserEnum(
			nil, // local
			2,   // level 2, USER_INFO_2
			netmgmtInternalApi.FILTER_NORMAL_ACCOUNT,
			&buf,
			netmgmtInternalApi.MAX_PREFERRED_LENGTH,
			&entriesread,
			&totalentries,
			&resumehandle,
		)
		if ret == netmgmtInternalApi.NERR_Success || ret == win.ERROR_MORE_DATA {
			var bufPtr = unsafe.Pointer(buf)
			for i := uint32(0); i < entriesread; i++ {
				var bufData = (*netmgmtInternalApi.USER_INFO_2)(unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample)))

				var user = &netmgmtModelsPb.User{
					Username:  encode.UTF16PtrToString(bufData.Usri2_name),
					Privilege: bufData.Usri2_priv,
					Flags:     bufData.Usri2_flags,
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
		default:
			return nil, status.Errorf(codes.Unknown, "Failed to get users (error: 0x%x)", ret)
		}
	}

	netmgmtInternalApi.NetApiBufferFree(buf)
	buf = nil

	return users, nil
}

func GetUser(user *netmgmtModelsPb.User_1) (*netmgmtModelsPb.User, error) {
	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	username, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}

	var buf = new(byte)

	ret, _, _ := netmgmtInternalApi.NetUserGetInfo(
		nil, // local
		username,
		2, // level 2, USER_INFO_2
		&buf,
	)

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return nil, status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_UserNotFound:
			return nil, status.Errorf(codes.NotFound, "User not found")
		default:
			return nil, status.Errorf(codes.Unknown, "Failed to get user (error: 0x%x)", ret)
		}
	}

	var bufPtr = unsafe.Pointer(buf)
	var bufData = (*netmgmtInternalApi.USER_INFO_2)(bufPtr)

	fetchedUser := &netmgmtModelsPb.User{
		Username:  encode.UTF16PtrToString(bufData.Usri2_name),
		Privilege: bufData.Usri2_priv,
		Flags:     bufData.Usri2_flags,
	}

	netmgmtInternalApi.NetApiBufferFree(buf)
	buf = nil

	return fetchedUser, nil
}

func GetUserLocalGroups(user *netmgmtModelsPb.User_1) ([]*netmgmtModelsPb.LocalGroup_1, error) {
	var localGroups []*netmgmtModelsPb.LocalGroup_1

	if user == nil {
		return nil, status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	username, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}

	var buf = new(byte)
	var entriesread uint32
	var totalentries uint32

	var bufDataSample netmgmtInternalApi.LOCALGROUP_USERS_INFO_0

	var ret uintptr

	for {
		ret, _, _ = netmgmtInternalApi.NetUserGetLocalGroups(
			nil, // local
			username,
			0, // level 0, LOCALGROUP_USERS_INFO_0
			0, // flags 0, none,
			&buf,
			netmgmtInternalApi.MAX_PREFERRED_LENGTH,
			&entriesread,
			&totalentries,
		)
		if ret == netmgmtInternalApi.NERR_Success || ret == win.ERROR_MORE_DATA {
			var bufPtr = unsafe.Pointer(buf)
			for i := uint32(0); i < entriesread; i++ {
				var bufData = (*netmgmtInternalApi.LOCALGROUP_USERS_INFO_0)(unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample)))

				var localGroup = &netmgmtModelsPb.LocalGroup_1{
					Groupname: encode.UTF16PtrToString(bufData.Lgrui0_name),
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
		case netmgmtInternalApi.NERR_BadUsername:
			return nil, status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_UserNotFound:
			return nil, status.Errorf(codes.NotFound, "User not found")
		default:
			return nil, status.Errorf(codes.Unknown, "Failed to get user local groups (error: 0x%x)", ret)
		}
	}

	netmgmtInternalApi.NetApiBufferFree(buf)
	buf = nil

	return localGroups, nil
}

func ChangeUserPassword(user *netmgmtModelsPb.User_3) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	username, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}
	usri1003_password, err := encode.UTF16PtrFromString(user.Password)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode password to UTF16")
	}

	var bufData = &netmgmtInternalApi.USER_INFO_1003{
		Usri1003_password: usri1003_password,
	}
	var bufPtr = unsafe.Pointer(bufData)
	var buf = (*byte)(bufPtr)

	var parm_err uint32

	ret, _, _ := netmgmtInternalApi.NetUserSetInfo(
		nil, // local
		username,
		1003, // level 1003, USER_INFO_1003
		buf,
		&parm_err,
	)

	user.Password = ""
	secure.ZeroMemoryUint16FromPtr(usri1003_password)
	usri1003_password = nil

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_BadPassword:
			return status.Errorf(codes.InvalidArgument, "Bad password")
		case netmgmtInternalApi.NERR_UserNotFound:
			return status.Errorf(codes.NotFound, "User not found")
		case netmgmtInternalApi.NERR_PasswordTooShort:
			return status.Errorf(codes.InvalidArgument, "Password does not meet the password policy requirements")
		case netmgmtInternalApi.NERR_LastAdmin:
			return status.Errorf(codes.FailedPrecondition, "Operation not allowed on the last administrative account")
		default:
			return status.Errorf(codes.Unknown, "Failed to change user password (error: 0x%x)", ret)
		}
	}

	return nil
}

func getUserFlags(user *netmgmtModelsPb.User_1) (uint32, error) {
	fetchedUser, err := GetUser(user)
	if err != nil {
		return 0, err
	}
	return fetchedUser.Flags, nil
}

func setUserFlags(user *netmgmtModelsPb.User_1, flags uint32) error {
	username, err := encode.UTF16PtrFromString(user.Username)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Unable to encode username to UTF16")
	}

	var bufData = &netmgmtInternalApi.USER_INFO_1008{
		Usri1008_flags: flags,
	}
	var bufPtr = unsafe.Pointer(bufData)
	var buf = (*byte)(bufPtr)

	var parm_err uint32

	ret, _, _ := netmgmtInternalApi.NetUserSetInfo(
		nil, // local
		username,
		1008, // level 1008, USER_INFO_1008
		buf,
		&parm_err,
	)

	if ret != netmgmtInternalApi.NERR_Success {
		switch ret {
		case netmgmtInternalApi.NERR_BadUsername:
			return status.Errorf(codes.InvalidArgument, "Bad username")
		case netmgmtInternalApi.NERR_UserNotFound:
			return status.Errorf(codes.NotFound, "User not found")
		case netmgmtInternalApi.NERR_LastAdmin:
			return status.Errorf(codes.FailedPrecondition, "Operation not allowed on the last administrative account")
		default:
			return status.Errorf(codes.Unknown, "Failed to set user flags (error: 0x%x)", ret)
		}
	}

	return nil
}

func EnableUser(user *netmgmtModelsPb.User_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	flags, err := getUserFlags(user)
	if err != nil {
		return err
	}

	flags &^= netmgmtInternalApi.UF_ACCOUNTDISABLE

	return setUserFlags(user, flags)
}

func DisableUser(user *netmgmtModelsPb.User_1) error {
	if user == nil {
		return status.Errorf(codes.InvalidArgument, "User cannot be nil")
	}

	flags, err := getUserFlags(user)
	if err != nil {
		return err
	}

	flags |= netmgmtInternalApi.UF_ACCOUNTDISABLE

	return setUserFlags(user, flags)
}
