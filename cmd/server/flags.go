package main

import "flag"

type flags struct {
	serverAddress string
}

func parseFlags() flags {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")

	flag.Parse()

	return flags{
		serverAddress: *serverAddress,
	}
}
