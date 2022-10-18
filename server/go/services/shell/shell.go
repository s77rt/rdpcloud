//go:build windows && amd64

package shell

import (
	shellServicePb "github.com/s77rt/rdpcloud/proto/go/services/shell"
)

type Server struct {
	shellServicePb.UnimplementedShellServer
}
