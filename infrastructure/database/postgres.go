package database

import (
	_ "github.com/lib/pq"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const driverName = "postgres"

func NewPostgresClient(configuration environment.Postgres) (*Provider, error) {
	connString := configuration.BuildConnectionString()
	gormConfig := &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: tablePrefix}}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{DSN: connString}), gormConfig)
	if err != nil {
		return nil, err
	}
	return &Provider{
		Conn: gormDB,
	}, nil
}
