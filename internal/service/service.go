package service

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/andidu/metrics/internal/handler"
	"github.com/andidu/metrics/internal/utils"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func NewMemStorage(db *sql.DB) handler.MemStorage {
	return &memStorageImpl{
		counters: make(map[string]int64),
		gauges:   make(map[string]float64),
		db:       db,
	}
}

type memStorageImpl struct {
	m        sync.Mutex
	counters map[string]int64
	gauges   map[string]float64
	db       *sql.DB
}

var errWrongInputParam = errors.New("wrong input params")

func (m *memStorageImpl) Gauges() (map[string]float64, error) {
	if m.db != nil {
		return tryRunDBOperationWithValue(
			func() (map[string]float64, error) {
				return dbGauges(m.db)
			},
		)
	}

	m.m.Lock()
	defer m.m.Unlock()

	gauges := make(map[string]float64, len(m.gauges))
	for k, v := range m.gauges {
		gauges[k] = v
	}
	return gauges, nil
}

func dbGauges(db *sql.DB) (map[string]float64, error) {
	rows, err := db.Query(
		"SELECT name, fvalue FROM metrics WHERE fvalue IS NOT NULL",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gauges := make(map[string]float64)

	for rows.Next() {
		var name string
		var value float64

		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		gauges[name] = value
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return gauges, nil
}

func (m *memStorageImpl) Counters() (map[string]int64, error) {
	if m.db != nil {
		return tryRunDBOperationWithValue(
			func() (map[string]int64, error) {
				return dbCounters(m.db)
			},
		)
	}

	m.m.Lock()
	defer m.m.Unlock()

	counters := make(map[string]int64, len(m.counters))
	for k, v := range m.counters {
		counters[k] = v
	}
	return counters, nil
}

func dbCounters(db *sql.DB) (map[string]int64, error) {
	rows, err := db.Query(
		"SELECT name, ivalue FROM metrics WHERE ivalue IS NOT NULL",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counters := make(map[string]int64)

	for rows.Next() {
		var name string
		var value int64

		err := rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}

		counters[name] = value
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return counters, nil
}

func (m *memStorageImpl) UpdateGauge(name string, value float64) error {
	if m.db != nil {
		return tryRunDBOperation(
			func() error {
				_, err := m.db.Exec(`INSERT INTO metrics (name, fvalue) VALUES ($1, $2) 
		ON CONFLICT (name) DO UPDATE SET fvalue = EXCLUDED.fvalue;`, name, value)
				return err
			},
		)
	}

	m.m.Lock()
	m.gauges[name] = value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) UpdateCounter(name string, value int64) error {
	if m.db != nil {
		return tryRunDBOperation(
			func() error {
				_, err := m.db.Exec(`INSERT INTO metrics (name, ivalue) VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET ivalue = COALESCE(metrics.ivalue, 0) + EXCLUDED.ivalue;`, name, value)
				return err
			},
		)
	}

	m.m.Lock()
	m.counters[name] += value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) OverrideCounter(name string, value int64) error {
	if m.db != nil {
		return tryRunDBOperation(
			func() error {
				_, err := m.db.Exec(`INSERT INTO metrics (name, ivalue) VALUES ($1, $2) 
		ON CONFLICT (name) DO UPDATE SET ivalue = EXCLUDED.ivalue;`, name, value)
				return err
			},
		)
	}

	m.m.Lock()
	m.counters[name] = value
	m.m.Unlock()
	return nil
}

func (m *memStorageImpl) Ping() bool {
	if m.db != nil {
		return m.db.Ping() == nil
	}
	return false
}

func (m *memStorageImpl) UpdateMetrics(gauges map[string]float64, counters map[string]int64) error {
	if m.db != nil {
		return tryRunDBOperation(
			func() error {
				return dbUpdateMetrics(m.db, gauges, counters)
			},
		)
	}

	m.m.Lock()
	defer m.m.Unlock()
	for name, value := range counters {
		m.counters[name] += value
	}
	for name, value := range gauges {
		m.gauges[name] = value
	}
	return nil
}

func dbUpdateMetrics(db *sql.DB, gauges map[string]float64, counters map[string]int64) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	for name, value := range gauges {
		_, err = tx.Exec(`INSERT INTO metrics (name, fvalue) VALUES ($1, $2) 
		ON CONFLICT (name) DO UPDATE SET fvalue = EXCLUDED.fvalue`, name, value)
		if err != nil {
			return err
		}
	}

	for name, value := range counters {
		_, err = tx.Exec(`INSERT INTO metrics (name, ivalue) VALUES ($1, $2)
        ON CONFLICT (name) DO UPDATE SET ivalue = COALESCE(metrics.ivalue, 0) + EXCLUDED.ivalue`, name, value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func tryRunDBOperation(operation func() error) error {
	_, err := tryRunDBOperationWithValue(
		func() (struct{}, error) {
			return struct{}{}, operation()
		},
	)
	return err
}

func tryRunDBOperationWithValue[T any](operation func() (T, error)) (T, error) {
	return utils.TryValue(
		operation,
		func(err error) bool {
			var pgErr *pgconn.PgError

			return errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code)
		},
	)
}
