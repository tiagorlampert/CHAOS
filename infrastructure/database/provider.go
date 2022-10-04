package database

import (
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"gorm.io/gorm"
	"log"
)

type Provider struct {
	Conn *gorm.DB
}

func NewProvider(configuration environment.Database) (*Provider, error) {
	switch {
	case configuration.Sqlite.IsValid():
		log.Println("Starting SQLite database")
		return NewSqliteClient(configuration.Sqlite)
	case configuration.Postgres.IsValid():
		log.Println("Starting PostgreSQL database")
		return NewPostgresClient(configuration.Postgres)
	default:
		return nil, internal.ErrNoDatabaseProvided
	}
}

func (p *Provider) Migrate() {
	p.Conn.AutoMigrate(
		&entities.User{},
		&entities.Device{},
		&entities.Auth{},
	)
}
