package bt

import (
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/jackpal/bencode-go"
)

type FindPeersParams struct {
	InfoHash   string
	PeerId     string
	Port       uint16
	Uploaded   int
	Downloaded int
	Left       int
}

func FindPeersHTTP(announce string, params FindPeersParams) ([]Peer, error) {
	base, err := url.Parse(announce)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Set("info_hash", params.InfoHash)
	vals.Set("peer_id", params.PeerId)
	vals.Set("port", fmt.Sprint(params.Port))
	vals.Set("uploaded", fmt.Sprint(params.Uploaded))
	vals.Set("downloaded", fmt.Sprint(params.Downloaded))
	vals.Set("left", fmt.Sprint(params.Left))

	base.RawQuery = vals.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dto := struct {
		Interval uint64
		Peers    []struct {
			Ip   string
			Port uint64
		}
		FailureReason string `bencode:"failure reason"`
	}{}
	err = bencode.Unmarshal(resp.Body, &dto)
	if err != nil {
		return nil, err
	}

	if dto.FailureReason != "" {
		return nil, fmt.Errorf("failed to get peers: %s", dto.FailureReason)
	}

	peers := make([]Peer, len(dto.Peers))
	for i, p := range dto.Peers {
		peers[i] = Peer{
			Ip:   net.IP(p.Ip),
			Port: uint16(p.Port),
		}
	}

	return peers, nil
}
