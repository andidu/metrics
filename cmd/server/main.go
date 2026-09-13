package main

import (
	"net/http"
	"time"

	config "github.com/andidu/metrics/internal/config/server"
	"github.com/andidu/metrics/internal/gzip"
	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/mfiles"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/andidu/metrics/internal/utils"
)

func main() {
	utils.InitLogger()
	defer utils.Logger.Sync()

	var flags, err = config.ParseConfig()
	if err != nil {
		utils.Logger.Errorln(err)
		return
	}

	var storage handler.MemStorage = restoreMemStorageOrCreateNew(flags.Store)
	if storage == nil {
		return
	}

	var h handler.Handler
	if flags.Store.StoreInterval != 0 {
		startPeriodicMemStorageBumps(storage, flags.Store.StoreInterval, flags.Store.FileStoragePath)
		h = handler.New(storage)
	} else {
		h = handler.New(
			mfiles.SavingStorage{
				MemStorage: storage,
				Filename:   flags.Store.FileStoragePath,
			},
		)
	}

	r := router.MetricsRouter(h)

	err = http.ListenAndServe(flags.ServerAddress, gzip.WithGzip(utils.WithLogging(r)))
	if err != nil {
		println("Server didn't start", err.Error())
	}
}

// returns MemStorage or nil on error
func restoreMemStorageOrCreateNew(storeFlags config.Store) handler.MemStorage {
	if storeFlags.Restore {
		storage, err := mfiles.Restore(storeFlags.FileStoragePath)
		if err != nil {
			utils.Logger.Errorln("Unable to restore storage state", err)
			return nil
		} else {
			return storage
		}
	} else {
		return service.NewMemStorage()
	}
}

func startPeriodicMemStorageBumps(m handler.MemStorage, interval uint, filename string) {
	go func() {
		for {
			time.Sleep(time.Duration(interval) * time.Second)

			err := mfiles.Save(filename, m)
			if err != nil {
				utils.Logger.Errorln("unable to save mem storage", err)
			}
		}
	}()
}
