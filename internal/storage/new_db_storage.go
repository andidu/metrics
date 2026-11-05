package storage

import (
	"database/sql"
	"errors"

	"github.com/andidu/metrics/internal/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

type DBStorage struct {
	DB *sql.DB
	m  *migrate.Migrate
}

func (d DBStorage) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.m != nil {
		d.m.Close()
	}
}

func NewDBStorage(dsn string) (DBStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return DBStorage{}, err
	}

	m, err := runDBMigrations(db)
	if err != nil {
		utils.Logger.Errorln(err)
		return DBStorage{db, m}, err
	}

	return DBStorage{db, m}, nil
}

func runDBMigrations(db *sql.DB) (*migrate.Migrate, error) {
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
