package database

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectDB(ctx context.Context, connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxIdleTime(10 * 60)
	db.SetConnMaxLifetime(30 * 60)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	
	return db
}
