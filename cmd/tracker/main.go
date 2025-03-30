package main

import (
	"github.com/Sp92535/internal/tracker"
)

func main() {
	// initialize exchange
	tracker.ExchangeInit()
	defer tracker.StopExchange()

	// run tracker server
	tracker.RunTrackerServer()
}
