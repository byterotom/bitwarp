package main

import (
	"github.com/byterotom/internal/tracker"
)

func main() {
	tracker.ExchangeInit()
	// run tracker server
	tracker.RunTrackerServer()
}
