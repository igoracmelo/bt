package main

import (
	"net"
	"testing"
)

func Test_Peer_Address(t *testing.T) {
	p := Peer{
		Ip:   net.IP("123.456.789.0"),
		Port: 54321,
	}

	want := "123.456.789.0:54321"
	got := p.Address()
	if got != want {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
