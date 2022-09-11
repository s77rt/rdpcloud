//go:build windows && amd64

package secauthn

import (
	secauthnServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthn"
)

type Server struct {
	secauthnServicePb.UnimplementedSecauthnServer
}
