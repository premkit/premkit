package main

import (
	"math/rand"
	"time"

	"github.com/premkit/premkit/daemon"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	// Start the HTTP server.
	daemon.Run()

	<-make(chan int)
}
