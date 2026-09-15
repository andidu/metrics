package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
)

type config struct {
	ServerAddress string
	Metrics       metrics
}

type metrics struct {
	RepeatInterval int
	PollInterval   int
}

var NoEnvVariableFound = errors.New("no env variable found")

func ParseConfig() (config, error) {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
	var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")

	flag.Parse()

	addressEnv, foundAddress := os.LookupEnv("ADDRESS")
	if foundAddress {
		serverAddress = &addressEnv
	}

	reportIntervalEnv, err := lookupEnvInt("REPORT_INTERVAL")
	if err != nil {
		return config{}, err
	}
	repeatInterval = &reportIntervalEnv

	pollIntervalEnv, err := lookupEnvInt("POLL_INTERVAL")
	if err != nil {
		return config{}, err
	}
	pollInterval = &pollIntervalEnv

	return config{
		ServerAddress: *serverAddress,
		Metrics: metrics{
			RepeatInterval: *repeatInterval,
			PollInterval:   *pollInterval,
		},
	}, nil
}

func lookupEnvInt(key string) (int, error) {
	str, found := os.LookupEnv(key)

	if !found {
		return 0, NoEnvVariableFound
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return value, nil
}
