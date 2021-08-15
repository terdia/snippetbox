package datasource

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"
	dbHost     = "DB_HOST"
	dbSchema   = "DB_SCHEMA"
)

var (
	db       *sql.DB
	username = os.Getenv(dbUsername)
	password = os.Getenv(dbPassword)
	host     = os.Getenv(dbHost)
	schema   = os.Getenv(dbSchema)
)

type connectionPool struct {
	*sql.DB
}

func NewConnectionPool() (*connectionPool, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, host, schema)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &connectionPool{db}, nil
}
