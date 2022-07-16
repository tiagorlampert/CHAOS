package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
)

type authSqliteRepository struct {
	dbClient *gorm.DB
}

func NewAuthRepository(dbClient *gorm.DB) repositories.Auth {
	return &authSqliteRepository{dbClient: dbClient}
}

func (s authSqliteRepository) Insert(auth entities.Auth) error {
	return s.dbClient.Create(&auth).Error
}

func (s authSqliteRepository) Update(auth entities.Auth) error {
	return s.dbClient.Model(&auth).Update(&auth).Error
}

func (s authSqliteRepository) GetFirst() (entities.Auth, error) {
	var auth entities.Auth
	if err := s.dbClient.Find(&auth).Error; err != nil {
		return entities.Auth{}, handleError(err)
	}
	return auth, nil
}
