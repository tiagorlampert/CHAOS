package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/tiagorlampert/CHAOS/entities"
	"os"
	"strings"
)

const (
	dialect     = `sqlite3`
	dbExtension = `.db`
)

type Database struct {
	Conn *gorm.DB
}

func NewSQLiteClient(dir, dbName string) (*Database, error) {
	dir = strings.TrimSuffix(dir, string(os.PathSeparator))
	dbConn, err := gorm.Open(dialect, fmt.Sprint(dir, string(os.PathSeparator), dbName, dbExtension))
	if err != nil {
		return nil, err
	}
	conn := &Database{Conn: dbConn}
	conn.Migrate()
	return conn, nil
}

func (d *Database) Migrate() {
	d.Conn.AutoMigrate(
		&entities.User{},
		&entities.Device{},
		&entities.System{},
	)
}
