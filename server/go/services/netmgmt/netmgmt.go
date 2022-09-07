//go:build windows && amd64

package netmgmt

import (
	netmgmtServicePb "github.com/s77rt/rdpcloud/proto/go/services/netmgmt"
)

type Server struct {
	netmgmtServicePb.UnimplementedNetmgmtServer
}
