package main

import "net"

type Peer struct {
	Ip   net.IP
	Port uint16
}
