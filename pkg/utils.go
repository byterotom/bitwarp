package pkg

import (
	"log"
	"net"
	"time"
)

// function to get local ip -> note: this is not public ip but the private ip in your network
func GetLocalIp() string {
	// dial google dns
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("error getting local ip: %v", err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String()
}

// function to calculate round trip time of a function
func RTT(f func()) float64 {
	start := time.Now()
	f()
	return time.Since(start).Seconds()
}
