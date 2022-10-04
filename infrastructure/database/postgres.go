package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const driverName = "postgres"

func NewPostgresClient(configuration environment.Postgres) (*Provider, error) {
	db, err := newConnection(configuration.BuildConnectionString())
	if err != nil {
		return nil, err
	}
	gormConfig := &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: tablePrefix}}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), gormConfig)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Conn: gormDB,
	}, nil
}

func newConnection(connString string) (*sql.DB, error) {
	db, err := sql.Open(driverName, connString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
