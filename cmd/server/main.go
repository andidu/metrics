package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	config "github.com/andidu/metrics/internal/config/server"
	"github.com/andidu/metrics/internal/gzip"
	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/mfiles"
	"github.com/andidu/metrics/internal/router"
	"github.com/andidu/metrics/internal/service"
	"github.com/andidu/metrics/internal/utils"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	utils.InitLogger()
	defer utils.Logger.Sync()

	var flags, err = config.ParseConfig()
	if err != nil {
		utils.Logger.Errorln(err)
		return
	}

	var db *sql.DB
	if flags.DB.Dsn == "" {
		db = nil
	} else {
		db, err = sql.Open("pgx", flags.DB.Dsn)
		if err != nil {
			utils.Logger.Errorln(err)
			return
		}
		defer db.Close()

		m, err := runDbMigrations(db)
		if m != nil {
			defer m.Close()
		}
		if err != nil {
			utils.Logger.Errorln(err)
			return
		}
	}
	var storage = restoreMemStorageOrCreateNew(flags.Store, db)
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
func restoreMemStorageOrCreateNew(storeFlags config.Store, db *sql.DB) handler.MemStorage {
	if storeFlags.Restore {
		storage, err := mfiles.Restore(storeFlags.FileStoragePath, db)
		if err != nil {
			utils.Logger.Errorln("Unable to restore storage state", err)
			return nil
		} else {
			return storage
		}
	} else {
		return service.NewMemStorage(db)
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

func runDbMigrations(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"pgx5",
		driver,
	)
	if err != nil {
		return m, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return m, err
	}
	return m, nil
}
