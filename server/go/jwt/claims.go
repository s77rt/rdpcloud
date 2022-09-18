//go:build windows && amd64

package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Privilege         uint32 `json:"privilege"`
	jwt.RegisteredClaims
}
