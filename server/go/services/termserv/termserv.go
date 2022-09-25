//go:build windows && amd64

package termserv

import (
	termservServicePb "github.com/s77rt/rdpcloud/proto/go/services/termserv"
)

type Server struct {
	termservServicePb.UnimplementedTermservServer
}
