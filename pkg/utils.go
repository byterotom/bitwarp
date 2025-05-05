package pkg

import (
	"log"
	"math/rand/v2"
	"net"
	"slices"
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

// function to get random numbers from a slice
func GetRandom(slice []uint64, n uint64) []uint64 {
	
	seed1 := uint64(time.Now().UnixNano())
    seed2 := seed1 ^ (seed1 >> 32)
    r := rand.New(rand.NewPCG(seed1, seed2))

	copied := slices.Clone(slice)
	r.Shuffle(len(copied), func(i, j int) {
		copied[i], copied[j] = copied[j], copied[i]
	})
	return copied[:min(uint64(len(copied)), n)]
}

// funtion to remove a number from sorted slice
func Remove(sorted []uint64, target uint64) []uint64 {
	l := 0
	r := len(sorted) - 1
	idx := -1
	for l <= r {
		mid := l + (r-l)/2
		if sorted[mid] == target {
			idx = mid
			break
		}
		if sorted[mid] < target {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	if idx != -1 {
		return slices.Delete(sorted, idx, idx+1)
	}
	return sorted
}
