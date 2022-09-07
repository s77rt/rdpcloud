//go:build windows && amd64

package netmgmt

import (
	"fmt"
	"unsafe"

	netmgmtModelsPb "github.com/s77rt/rdpcloud/proto/go/models/netmgmt"
	"github.com/s77rt/rdpcloud/server/go/internal/win"
	"github.com/s77rt/rdpcloud/server/go/internal/win/headers"
	netmgmtInternalApi "github.com/s77rt/rdpcloud/server/go/internal/win/win32/netmgmt"
)

func GetUsers() ([]*netmgmtModelsPb.User, error) {
	var users []*netmgmtModelsPb.User

	var buf = new(byte)
	var entriesread uint32
	var totalentries uint32
	var resume_handle uint32

	var bufDataSample headers.USER_INFO_2

	var ret uintptr

	for {
		ret, _, _ = netmgmtInternalApi.NetUserEnum(
			nil, // local
			2,   // level 2, USER_INFO_2
			headers.FILTER_NORMAL_ACCOUNT,
			&buf,
			headers.MAX_PREFERRED_LENGTH,
			&entriesread,
			&totalentries,
			&resume_handle,
		)
		if ret == headers.NERR_Success || ret == headers.ERROR_MORE_DATA {
			var bufPtr = unsafe.Pointer(buf)
			for i := uint32(0); i < entriesread; i++ {
				bufPtr = unsafe.Pointer(uintptr(bufPtr) + uintptr(i)*unsafe.Sizeof(bufDataSample))
				var bufData = (*headers.USER_INFO_2)(bufPtr)

				var user = &netmgmtModelsPb.User{
					Username:  win.UTF16PtrToString(bufData.Usri2_name),
					Privilege: bufData.Usri2_priv,
					Flags:     bufData.Usri2_flags,
					FullName:  win.UTF16PtrToString(bufData.Usri2_full_name),
				}
				users = append(users, user)
			}

			if ret == headers.ERROR_MORE_DATA {
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

	netmgmtInternalApi.NetApiBufferFree(buf)

	if ret != headers.NERR_Success {
		return nil, fmt.Errorf("Failed to get local user accounts (error: %d)", ret)
	}

	return users, nil
}

func AddUser(user *netmgmtModelsPb.User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	usri2_name, err := win.UTF16PtrFromString(user.Username)
	if err != nil {
		return fmt.Errorf("unable to encode username to UTF16")
	}
	usri2_password, err := win.UTF16PtrFromString(user.Password)
	if err != nil {
		return fmt.Errorf("unable to encode password to UTF16")
	}
	usri2_full_name, err := win.UTF16PtrFromString(user.FullName)
	if err != nil {
		return fmt.Errorf("unable to encode full name to UTF16")
	}

	var bufData = &headers.USER_INFO_2{
		Usri2_name:      usri2_name,
		Usri2_password:  usri2_password,
		Usri2_priv:      headers.USER_PRIV_USER,
		Usri2_flags:     headers.UF_NORMAL_ACCOUNT | headers.UF_SCRIPT | headers.UF_DONT_EXPIRE_PASSWD,
		Usri2_full_name: usri2_full_name,
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

	if ret != headers.NERR_Success {
		return fmt.Errorf("Failed to add user (error: %d)", ret)
	}

	return nil
}
