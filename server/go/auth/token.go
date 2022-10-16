//go:build windows && amd64

package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/s77rt/rdpcloud/server/go/config"
)

func NewTokenWithUserClaims(userClaims UserClaims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwtUserClaims{
		PreferredUsername: userClaims.PreferredUsername,
		Privilege:         userClaims.Privilege,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userClaims.UserSID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenLifetime * time.Second)),
		},
	})
}

func TokenSignedString(token *jwt.Token) (string, error) {
	return token.SignedString(config.Secret)
}

func ParseTokenWithUserClaims(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &jwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect

		signingMethod, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		if signingMethod != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.Secret, nil
	})
}

func TokenUserClaims(token *jwt.Token) (*UserClaims, error) {
	tokenJwtUserClaims, ok := token.Claims.(*jwtUserClaims)
	if !ok {
		return nil, fmt.Errorf("Unexpected claims")
	}

	return &UserClaims{
		UserSID:           tokenJwtUserClaims.RegisteredClaims.Subject,
		PreferredUsername: tokenJwtUserClaims.PreferredUsername,
		Privilege:         tokenJwtUserClaims.Privilege,
	}, nil
}
