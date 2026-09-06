package config

import (
	"flag"
)

type config struct {
	ServerAddress string
	Metrics       metrics
}

type metrics struct {
	RepeatInterval int
	PollInterval   int
}

func ParseConfig() config {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
	var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")

	flag.Parse()

	return config{
		ServerAddress: *serverAddress,
		Metrics: metrics{
			RepeatInterval: *repeatInterval,
			PollInterval:   *pollInterval,
		},
	}
}
