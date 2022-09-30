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
// If an rpc method level is not specified, you are expected to deny access for all as a security fallback
var AceessLevel = map[string]uint32{
	"/services.fileio.Fileio/GetQuotaState":                  1,
	"/services.fileio.Fileio/SetQuotaState":                  1,
	"/services.fileio.Fileio/GetDefaultQuota":                1,
	"/services.fileio.Fileio/SetDefaultQuota":                1,
	"/services.fileio.Fileio/GetUsersQuotaEntries":           1,
	"/services.fileio.Fileio/GetUserQuotaEntry":              1,
	"/services.fileio.Fileio/SetUserQuotaEntry":              1,
	"/services.fileio.Fileio/DeleteUserQuotaEntry":           1,
	"/services.fileio.Fileio/GetVolumes":                     1,
	"/services.netmgmt.Netmgmt/AddUser":                      1,
	"/services.netmgmt.Netmgmt/DeleteUser":                   1,
	"/services.netmgmt.Netmgmt/GetUsers":                     1,
	"/services.netmgmt.Netmgmt/GetUser":                      1,
	"/services.netmgmt.Netmgmt/GetUserLocalGroups":           1,
	"/services.netmgmt.Netmgmt/ChangeUserPassword":           1,
	"/services.netmgmt.Netmgmt/EnableUser":                   1,
	"/services.netmgmt.Netmgmt/DisableUser":                  1,
	"/services.netmgmt.Netmgmt/AddUserToLocalGroup":          1,
	"/services.netmgmt.Netmgmt/RemoveUserFromLocalGroup":     1,
	"/services.netmgmt.Netmgmt/GetLocalGroups":               1,
	"/services.netmgmt.Netmgmt/GetUsersInLocalGroup":         1,
	"/services.secauthn.Secauthn/LogonUser":                  0,
	"/services.secauthz.Secauthz/LookupAccountSidByUsername": 1,
	"/services.secauthz.Secauthz/LookupAccountUsernameBySid": 1,
	"/services.termserv.Termserv/LogoffUser":                 1,
}
