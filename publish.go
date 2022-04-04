package main

import (
	"github.com/nats-io/stan.go"
)

func main() {
	sc, _ := stan.Connect("test-cluster", "client1")

	// Simple Synchronous Publisher
	sc.Publish("Lorem2", []byte("Mauris pretium purus odio, a luctus odio ornare molestie. Nam ac tellus vel dolor accumsan dictum et eu diam. Vestibulum et sapien sollicitudin, convallis odio sed, lobortis enim. ")) // does not return until an ack has been received from NATS Streaming

	// Simple Async Subscribe

	// Close connection
	sc.Close()

}
