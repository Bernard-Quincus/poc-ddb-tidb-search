package db

import "database/sql"

type DB interface {
	Get(any) (any, error)
	Put(...any) (any, error)
	Delete(any) error
	Search(any) (any, error)
	Close() error
	SetTableName(string) error
	GetTiDBConn() *sql.DB
}
