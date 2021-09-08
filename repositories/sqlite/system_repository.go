package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
)

type systemRepository struct {
	database *gorm.DB
}

func NewSystemRepository(database *gorm.DB) repositories.System {
	return &systemRepository{
		database: database,
	}
}

func (s systemRepository) Insert(config entities.System) error {
	return s.database.Create(&config).Error
}

func (s systemRepository) Update(config entities.System) (*entities.System, error) {
	if err := s.database.Model(&config).Update(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (s systemRepository) Get() (*entities.System, error) {
	var system entities.System
	if err := s.database.Find(&system).Error; err != nil {
		return nil, err
	}
	return &system, nil
}
