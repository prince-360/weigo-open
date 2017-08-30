package database

import (
	"github.com/jackc/pgx"
)

var (
	pool *pgx.ConnPool
)

// Connect .
func Connect(uri string) {
	conf, err := pgx.ParseURI(uri)
	if err != nil {
		panic(err)
	}
	confPool := pgx.ConnPoolConfig{
		ConnConfig:     conf,
		MaxConnections: 2,
	}
	tmpPool, err := pgx.NewConnPool(confPool)
	if err != nil {
		panic(err)
	}
	pool = tmpPool
}
