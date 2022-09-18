//go:build windows && amd64

package secauthz

import (
	secauthzServicePb "github.com/s77rt/rdpcloud/proto/go/services/secauthz"
)

type Server struct {
	secauthzServicePb.UnimplementedSecauthzServer
}
