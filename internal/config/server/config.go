package config

import (
	"flag"
	"os"
	"strconv"

	"github.com/andidu/metrics/internal/utils"
)

type config struct {
	ServerAddress string
	Store         store
}

type store struct {
	StoreInterval   uint
	FileStoragePath string
	Restore         bool
}

func ParseConfig() config {
	var serverAddress = flag.String("a", "localhost:8080", "Server IP addres")
	var storeInterval = flag.Uint("i", 300, "Interval between metrics disk dumps")
	var fileStoragePath = flag.String("f", "dumps/metrics.txt", "File path to the metrics dump file")
	var restore = flag.Bool("r", false, "Whether to restore metrics state from fileStoragePath")

	flag.Parse()

	addressEnv, found := os.LookupEnv("ADDRESS")
	if found {
		serverAddress = &addressEnv
	}
	storeIntervalEnv, found := os.LookupEnv("STORE_INTERVAL")
	if found {
		uint64Value, err := strconv.ParseUint(storeIntervalEnv, 10, 64)
		if err != nil {
			utils.Logger.Errorln("unable to parse STORE_INTERVAL value")
		} else {
			uintValue := uint(uint64Value)
			storeInterval = &uintValue
		}
	}
	fileStoragePathEnv, found := os.LookupEnv("FILE_STORAGE_PATH")
	if found {
		fileStoragePath = &fileStoragePathEnv
	}
	restoreEnv, found := os.LookupEnv("RESTORE")
	if found {
		boolValue, err := strconv.ParseBool(restoreEnv)
		if err != nil {
			utils.Logger.Errorln("unable to parse RESTORE value")
		} else {
			restore = &boolValue
		}
	}

	return config{
		ServerAddress: *serverAddress,
		Store: store{
			StoreInterval:   *storeInterval,
			FileStoragePath: *fileStoragePath,
			Restore:         *restore,
		},
	}
}
