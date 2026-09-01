package main

import (
	"github.com/spf13/pflag"
)

var serverAddress = pflag.String("a", "localhost:8080", "Server IP addres")
var repeatInterval = pflag.Int("r", 10, "Metrics push repeat interval in seconds")
var pollInterval = pflag.Int("p", 2, "Metrics collection repeat interval")
