package bt

import (
	"fmt"
	"net"
)

type Peer struct {
	Ip   net.IP
	Port uint16
}

func (p *Peer) Address() string {
	return net.JoinHostPort(string(p.Ip), fmt.Sprint(p.Port))
}
