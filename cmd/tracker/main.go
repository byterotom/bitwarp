package main

import (
	"github.com/Sp92535/internal/tracker"
)

func main() {
	tracker.ExchangeInit()
	// run tracker server
	tracker.RunTrackerServer()
}
