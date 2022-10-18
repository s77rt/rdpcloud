//go:build windows && amd64

package shutdown

import (
	shutdownServicePb "github.com/s77rt/rdpcloud/proto/go/services/shutdown"
)

type Server struct {
	shutdownServicePb.UnimplementedShutdownServer
}
