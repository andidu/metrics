package config

import (
	"flag"
	"os"
)

type config struct {
	ServerAddress string
}

func ParseConfig() config {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	flag.Parse()

	addressEnv, foundAdress := os.LookupEnv("ADDRESS")
	if foundAdress {
		serverAddress = &addressEnv
	}

	return config{
		ServerAddress: *serverAddress,
	}
}
