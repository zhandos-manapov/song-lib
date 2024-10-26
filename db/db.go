package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database interface {
	GetPool() *pgxpool.Pool
	Connect() (Database, error)
	Disconnect()
}

type DbConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

type database struct {
	pool    *pgxpool.Pool
	context context.Context
	config  DbConfig
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	return &database{
		context: ctx,
		config:  config,
	}
}

func (db *database) Connect() (Database, error) {
	psqlInfo := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		db.config.User,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.Name,
	)

	if pool, err := pgxpool.New(db.context, psqlInfo); err != nil {
		return nil, err
	} else {
		log.Println("Successfully connected to the database!")
		db.pool = pool
		return db, nil
	}
}

func (db *database) Disconnect() {
	db.pool.Close()
}

func (db *database) GetPool() *pgxpool.Pool {
	return db.pool
}
