package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tiagorlampert/CHAOS/entities"
)

const dialect = `sqlite3`

type database struct {
	Conn *gorm.DB
}

func NewSQLiteClient(dbName string) (*database, error) {
	db, err := gorm.Open(dialect, fmt.Sprint("database/", dbName, `.db`))
	if err != nil {
		return nil, err
	}
	conn := &database{Conn: db}
	conn.Migrate()
	return conn, nil
}

func (d *database) Migrate() {
	d.Conn.AutoMigrate(
		&entities.User{},
		&entities.Device{},
		&entities.System{},
	)
}
