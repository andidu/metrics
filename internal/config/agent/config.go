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

func ParseConfig() config {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var repeatInterval = flag.Int("r", 10, "Metrics push repeat interval in seconds")
	var pollInterval = flag.Int("p", 2, "Metrics collection repeat interval")

	flag.Parse()

	addressEnv, foundAddress := os.LookupEnv("ADDRESS")
	if foundAddress {
		serverAddress = &addressEnv
	}

	reportIntervalEnv, foundReportInterval := lookupEnvInt("REPORT_INTERVAL")
	if foundReportInterval {
		repeatInterval = &reportIntervalEnv
	}

	pollIntervalEnv, foundPullInterval := lookupEnvInt("POLL_INTERVAL")
	if foundPullInterval {
		pollInterval = &pollIntervalEnv
	}

	return config{
		ServerAddress: *serverAddress,
		Metrics: metrics{
			RepeatInterval: *repeatInterval,
			PollInterval:   *pollInterval,
		},
	}
}

func lookupEnvInt(key string) (int, bool) {
	str, found := os.LookupEnv(key)

	if !found {
		return 0, false
	}

	value, err := strconv.Atoi(str)
	if err != nil {
		return 0, false
	}

	return value, true
}
