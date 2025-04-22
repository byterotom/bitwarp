package pkg

import (
	"log"
	"net"
	"time"
)

func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("error getting local ip: %v", err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String()
}

func RTT[T any](f func()) float64 {
	start := time.Now()
	f()
	return time.Since(start).Seconds()
}
