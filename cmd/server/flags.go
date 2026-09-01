package main

import (
	"github.com/spf13/pflag"
)

var serverAddress = pflag.String("a", "localhost:8080", "Server IP addres")
