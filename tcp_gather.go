package main

import (
	"net"
	"time"
)

func tcpGather(ip string, ports []string) map[string]bool {
	// check emqx 1883, 8083 port

	results := make(map[string]bool)
	for _, port := range ports {
		address := net.JoinHostPort(ip, port)
		// 3 second timeout
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			results[port] = false
			// todo log handler
		} else {
			if conn != nil {
				results[port] = true
				_ = conn.Close()
			} else {
				results[port] = false
			}
		}
	}
	return results
}
