package config

import "flag"

type config struct {
	ServerAddress string
}

func ParseConfig() config {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")

	flag.Parse()

	return config{
		ServerAddress: *serverAddress,
	}
}
