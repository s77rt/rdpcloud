//go:build windows && amd64

package auth

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	UserSID           string
	PreferredUsername string
	Privilege         uint32
}

type jwtUserClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Privilege         uint32 `json:"privilege"`
	jwt.RegisteredClaims
}
