package main

import (
	ping "github.com/digineo/go-ping"
	"net"
	"time"
)

func doPing(host string) (rtt time.Duration, err error) {
	var remoteAddr *net.IPAddr
	var pinger *ping.Pinger

	if r, err := net.ResolveIPAddr("ip4", host); err != nil {
		return rtt, err
	} else {
		remoteAddr = r
	}

	if p, err := ping.New("0.0.0.0", ""); err != nil {
		return rtt, err
	} else {
		pinger = p
	}

	defer pinger.Close()

	if pinger.PayloadSize() != uint16(56) {
		pinger.SetPayloadSize(uint16(56))
	}

	return pinger.PingAttempts(remoteAddr, time.Second, 2)
}
