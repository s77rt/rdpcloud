//go:build windows && amd64

package netmgmt

import (
	"unsafe"
)

type USER_INFO_1003 struct {
	Usri1003_password *uint16
}

type USER_INFO_2 struct {
	Usri2_name           *uint16
	Usri2_password       *uint16
	Usri2_password_age   uint32
	Usri2_priv           uint32
	Usri2_home_dir       *uint16
	Usri2_comment        *uint16
	Usri2_flags          uint32
	Usri2_script_path    *uint16
	Usri2_auth_flags     uint32
	Usri2_full_name      *uint16
	Usri2_usr_comment    *uint16
	Usri2_parms          *uint16
	Usri2_workstations   *uint16
	Usri2_last_logon     uint32
	Usri2_last_logoff    uint32
	Usri2_acct_expires   uint32
	Usri2_max_storage    uint32
	Usri2_units_per_week uint32
	Usri2_logon_hours    *byte
	Usri2_bad_pw_count   uint32
	Usri2_num_logons     uint32
	Usri2_logon_server   *uint16
	Usri2_country_code   uint32
	Usri2_code_page      uint32
}

type LOCALGROUP_USERS_INFO_0 struct {
	Lgrui0_name *uint16
}

const UF_SCRIPT = 0x0001
const UF_ACCOUNTDISABLE = 0x0002
const UF_HOMEDIR_REQUIRED = 0x0008
const UF_LOCKOUT = 0x0010
const UF_PASSWD_NOTREQD = 0x0020
const UF_PASSWD_CANT_CHANGE = 0x0040
const UF_ENCRYPTED_TEXT_PASSWORD_ALLOWED = 0x0080

const UF_TEMP_DUPLICATE_ACCOUNT = 0x0100
const UF_NORMAL_ACCOUNT = 0x0200
const UF_INTERDOMAIN_TRUST_ACCOUNT = 0x0800
const UF_WORKSTATION_TRUST_ACCOUNT = 0x1000
const UF_SERVER_TRUST_ACCOUNT = 0x2000

const UF_MACHINE_ACCOUNT_MASK = (UF_INTERDOMAIN_TRUST_ACCOUNT | UF_WORKSTATION_TRUST_ACCOUNT | UF_SERVER_TRUST_ACCOUNT)
const UF_ACCOUNT_TYPE_MASK = (UF_TEMP_DUPLICATE_ACCOUNT | UF_NORMAL_ACCOUNT | UF_INTERDOMAIN_TRUST_ACCOUNT | UF_WORKSTATION_TRUST_ACCOUNT | UF_SERVER_TRUST_ACCOUNT)

const UF_DONT_EXPIRE_PASSWD = 0x10000
const UF_MNS_LOGON_ACCOUNT = 0x20000
const UF_SMARTCARD_REQUIRED = 0x40000
const UF_TRUSTED_FOR_DELEGATION = 0x80000
const UF_NOT_DELEGATED = 0x100000
const UF_USE_DES_KEY_ONLY = 0x200000
const UF_DONT_REQUIRE_PREAUTH = 0x400000
const UF_PASSWORD_EXPIRED = 0x800000
const UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION = 0x1000000
const UF_NO_AUTH_DATA_REQUIRED = 0x2000000

const UF_SETTABLE_BITS = (UF_SCRIPT | UF_ACCOUNTDISABLE | UF_LOCKOUT | UF_HOMEDIR_REQUIRED | UF_PASSWD_NOTREQD | UF_PASSWD_CANT_CHANGE | UF_ACCOUNT_TYPE_MASK | UF_DONT_EXPIRE_PASSWD | UF_MNS_LOGON_ACCOUNT | UF_ENCRYPTED_TEXT_PASSWORD_ALLOWED | UF_SMARTCARD_REQUIRED | UF_TRUSTED_FOR_DELEGATION | UF_NOT_DELEGATED | UF_USE_DES_KEY_ONLY | UF_DONT_REQUIRE_PREAUTH | UF_PASSWORD_EXPIRED | UF_TRUSTED_TO_AUTHENTICATE_FOR_DELEGATION | UF_NO_AUTH_DATA_REQUIRED)

const FILTER_TEMP_DUPLICATE_ACCOUNT = (0x0001)
const FILTER_NORMAL_ACCOUNT = (0x0002)

const FILTER_INTERDOMAIN_TRUST_ACCOUNT = (0x0008)
const FILTER_WORKSTATION_TRUST_ACCOUNT = (0x0010)
const FILTER_SERVER_TRUST_ACCOUNT = (0x0020)

const LG_INCLUDE_INDIRECT = (0x0001)

const USER_NAME_PARMNUM = 1
const USER_PASSWORD_PARMNUM = 3
const USER_PASSWORD_AGE_PARMNUM = 4
const USER_PRIV_PARMNUM = 5
const USER_HOME_DIR_PARMNUM = 6
const USER_COMMENT_PARMNUM = 7
const USER_FLAGS_PARMNUM = 8
const USER_SCRIPT_PATH_PARMNUM = 9
const USER_AUTH_FLAGS_PARMNUM = 10
const USER_FULL_NAME_PARMNUM = 11
const USER_USR_COMMENT_PARMNUM = 12
const USER_PARMS_PARMNUM = 13
const USER_WORKSTATIONS_PARMNUM = 14
const USER_LAST_LOGON_PARMNUM = 15
const USER_LAST_LOGOFF_PARMNUM = 16
const USER_ACCT_EXPIRES_PARMNUM = 17
const USER_MAX_STORAGE_PARMNUM = 18
const USER_UNITS_PER_WEEK_PARMNUM = 19
const USER_LOGON_HOURS_PARMNUM = 20
const USER_PAD_PW_COUNT_PARMNUM = 21
const USER_NUM_LOGONS_PARMNUM = 22
const USER_LOGON_SERVER_PARMNUM = 23
const USER_COUNTRY_CODE_PARMNUM = 24
const USER_CODE_PAGE_PARMNUM = 25
const USER_PRIMARY_GROUP_PARMNUM = 51
const USER_PROFILE = 52
const USER_PROFILE_PARMNUM = 52
const USER_HOME_DIR_DRIVE_PARMNUM = 53

const USER_PRIV_MASK = 0x3
const USER_PRIV_GUEST = 0
const USER_PRIV_USER = 1
const USER_PRIV_ADMIN = 2

var (
	procNetUserAdd            = modnetapi32.NewProc("NetUserAdd")
	procNetUserDel            = modnetapi32.NewProc("NetUserDel")
	procNetUserEnum           = modnetapi32.NewProc("NetUserEnum")
	procNetUserGetInfo        = modnetapi32.NewProc("NetUserGetInfo")
	procNetUserGetLocalGroups = modnetapi32.NewProc("NetUserGetLocalGroups")
	procNetUserSetInfo        = modnetapi32.NewProc("NetUserSetInfo")
)

func NetUserAdd(servername *uint16, level uint32, buf *byte, parm_err *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserAdd.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parm_err)),
	)
	return
}

func NetUserDel(servername *uint16, username *uint16) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserDel.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(username)),
	)
	return
}

func NetUserEnum(servername *uint16, level uint32, filter uint32, buf **byte, prefmaxlen uint32, entriesread *uint32, totalentries *uint32, resumehandle *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserEnum.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(level),
		uintptr(filter),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefmaxlen),
		uintptr(unsafe.Pointer(entriesread)),
		uintptr(unsafe.Pointer(totalentries)),
		uintptr(unsafe.Pointer(resumehandle)),
	)
	return
}

func NetUserGetInfo(servername *uint16, username *uint16, level uint32, buf **byte) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserGetInfo.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(username)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
	)
	return
}

func NetUserGetLocalGroups(servername *uint16, username *uint16, level uint32, flags uint32, buf **byte, prefmaxlen uint32, entriesread *uint32, totalentries *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserGetLocalGroups.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(username)),
		uintptr(level),
		uintptr(flags),
		uintptr(unsafe.Pointer(buf)),
		uintptr(prefmaxlen),
		uintptr(unsafe.Pointer(entriesread)),
		uintptr(unsafe.Pointer(totalentries)),
	)
	return
}

func NetUserSetInfo(servername *uint16, username *uint16, level uint32, buf *byte, parm_err *uint32) (r1, r2 uintptr, lastErr error) {
	r1, r2, lastErr = procNetUserSetInfo.Call(
		uintptr(unsafe.Pointer(servername)),
		uintptr(unsafe.Pointer(username)),
		uintptr(level),
		uintptr(unsafe.Pointer(buf)),
		uintptr(unsafe.Pointer(parm_err)),
	)
	return
}
