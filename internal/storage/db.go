package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(32)
	db.SetConnMaxLifetime(20 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("clickhouse ping failed: %w", err)
	}
	return db, nil
}
