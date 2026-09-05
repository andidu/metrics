package main

import (
	"flag"
)

type flags struct {
	serverAddress string
	metrics       metrics
}

type metrics struct {
	repeatInterval int
	pollInterval   int
}

func parseFlags() flags {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
	var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")

	flag.Parse()

	return flags{
		serverAddress: *serverAddress,
		metrics: metrics{
			repeatInterval: *repeatInterval,
			pollInterval:   *pollInterval,
		},
	}
}
