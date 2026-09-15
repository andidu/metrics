package config

import (
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

func ParseConfig() (config, error) {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
	var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")

	flag.Parse()

	addressEnv, foundAddress := os.LookupEnv("ADDRESS")
	if foundAddress {
		serverAddress = &addressEnv
	}

	reportIntervalEnv, found, err := lookupEnvInt("REPORT_INTERVAL")
	if err != nil {
		return config{}, err
	}
	if found {
		repeatInterval = &reportIntervalEnv
	}

	pollIntervalEnv, found, err := lookupEnvInt("POLL_INTERVAL")
	if err != nil {
		return config{}, err
	}
	if found {
		pollInterval = &pollIntervalEnv
	}

	return config{
		ServerAddress: *serverAddress,
		Metrics: metrics{
			RepeatInterval: *repeatInterval,
			PollInterval:   *pollInterval,
		},
	}, nil
}

// return the int value, a flag whether it was found
// and error in case it was found and there was a parsing error
func lookupEnvInt(key string) (int, bool, error) {
	str, found := os.LookupEnv(key)

	if !found {
		return 0, false, nil
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, true, err
	}

	return value, true, nil
}
