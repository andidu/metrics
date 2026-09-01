package main

import (
	"flag"
)

var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")
