//go:build windows && amd64

package auth

import (
	"context"
	"fmt"
)

type contextKey int

const userClaimsKey contextKey = 0

func NewContextWithUserClaims(ctx context.Context, userClaims *UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, userClaims)
}

func UserClaimsFromContext(ctx context.Context) (*UserClaims, error) {
	contextUserClaims := ctx.Value(userClaimsKey)
	if contextUserClaims == nil {
		return nil, fmt.Errorf("Missing claims")
	}

	userClaims, ok := contextUserClaims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Unexpected claims")
	}

	return userClaims, nil
}
