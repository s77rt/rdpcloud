//go:build windows && amd64

package config

import "crypto/rand"

func init() {
	Secret = make([]byte, 64)
	_, err := rand.Read(Secret)
	if err != nil {
		panic(err)
	}
}

var (
	Secret []byte
)

const TokenLifetime = 60 * 60 // in seconds

// AceessLevel is a map where the key is the full rpc method
// and the value is the minimum required level
// the required level is an analogy of the Windows User Privilege (0, 1, 2 for GUEST, USER, ADMIN resp)
// 0 => All
// 1 => Users
// 2 => Admins
// If a rpc method does not have a defined access level, you are expected to deny access for all as a security fallback
var AceessLevel = map[string]uint32{
	"/services.secauthn.Secauthn/LogonUser": 0,

	"/services.fileio.Fileio/GetMyUserQuotaEntry":    1,
	"/services.netmgmt.Netmgmt/ChangeMyUserPassword": 1,
	"/services.netmgmt.Netmgmt/GetMyUser":            1,
	"/services.netmgmt.Netmgmt/GetMyUserLocalGroups": 1,
	"/services.sysinfo.Sysinfo/GetUptime":            1,
	"/services.termserv.Termserv/LogoffMyUser":       1,

	"/services.fileio.Fileio/DeleteUserQuotaEntry":           2,
	"/services.fileio.Fileio/GetDefaultQuota":                2,
	"/services.fileio.Fileio/GetQuotaState":                  2,
	"/services.fileio.Fileio/GetUserQuotaEntry":              2,
	"/services.fileio.Fileio/GetUsersQuotaEntries":           2,
	"/services.fileio.Fileio/GetVolumes":                     2,
	"/services.fileio.Fileio/SetDefaultQuota":                2,
	"/services.fileio.Fileio/SetQuotaState":                  2,
	"/services.fileio.Fileio/SetUserQuotaEntry":              2,
	"/services.netmgmt.Netmgmt/AddUser":                      2,
	"/services.netmgmt.Netmgmt/AddUserToLocalGroup":          2,
	"/services.netmgmt.Netmgmt/ChangeUserPassword":           2,
	"/services.netmgmt.Netmgmt/DeleteUser":                   2,
	"/services.netmgmt.Netmgmt/DisableUser":                  2,
	"/services.netmgmt.Netmgmt/EnableUser":                   2,
	"/services.netmgmt.Netmgmt/GetLocalGroups":               2,
	"/services.netmgmt.Netmgmt/GetUser":                      2,
	"/services.netmgmt.Netmgmt/GetUserLocalGroups":           2,
	"/services.netmgmt.Netmgmt/GetUsers":                     2,
	"/services.netmgmt.Netmgmt/GetUsersInLocalGroup":         2,
	"/services.netmgmt.Netmgmt/RemoveUserFromLocalGroup":     2,
	"/services.secauthz.Secauthz/LookupAccountSidByUsername": 2,
	"/services.secauthz.Secauthz/LookupAccountUsernameBySid": 2,
	"/services.shell.Shell/DeleteProfile":                    2,
	"/services.termserv.Termserv/LogoffUser":                 2,
}
