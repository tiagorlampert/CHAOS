package database

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/internal/environment"
	"github.com/tiagorlampert/CHAOS/internal/utils/constants"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strings"
)

const (
	dbExtension = `.db`
)

func NewSqliteClient(configuration environment.Sqlite) (*Provider, error) {
	dir := strings.TrimSuffix(constants.DatabaseDirectory, string(os.PathSeparator))
	gormConfig := &gorm.Config{NamingStrategy: schema.NamingStrategy{TablePrefix: tablePrefix}}
	gormDB, err := gorm.Open(sqlite.Open(fmt.Sprint(dir, string(os.PathSeparator), configuration.DatabaseName, dbExtension)), gormConfig)
	if err != nil {
		return nil, err
	}
	return &Provider{Conn: gormDB}, nil
}
