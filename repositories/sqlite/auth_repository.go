package sqlite

import (
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"gorm.io/gorm"
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

func (s authSqliteRepository) Update(auth *entities.Auth) error {
	return s.dbClient.Model(&auth).Updates(&auth).Error
}

func (s authSqliteRepository) GetFirst() (*entities.Auth, error) {
	var auth entities.Auth
	tx := s.dbClient.Find(&auth)
	if tx.Error != nil {
		return nil, handleError(tx.Error)
	}
	if tx.RowsAffected == 0 {
		return nil, repositories.ErrNotFound
	}
	return &auth, nil
}
