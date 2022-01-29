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

type SqliteClient struct {
	Conn *gorm.DB
}

func NewSqliteClient(dir, dbName string) (*SqliteClient, error) {
	dir = strings.TrimSuffix(dir, string(os.PathSeparator))
	dbConn, err := gorm.Open(dialect, fmt.Sprint(dir, string(os.PathSeparator), dbName, dbExtension))
	if err != nil {
		return nil, err
	}
	conn := &SqliteClient{Conn: dbConn}
	conn.Migrate()
	return conn, nil
}

func (d *SqliteClient) Migrate() {
	d.Conn.AutoMigrate(
		&entities.User{},
		&entities.Device{},
		&entities.Auth{},
	)
}
