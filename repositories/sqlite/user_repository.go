package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
)

type userRepository struct {
	database *gorm.DB
}

func NewUserRepository(database *gorm.DB) repositories.User {
	return &userRepository{database: database}
}

func (u userRepository) Insert(user entities.User) error {
	return u.database.Create(&user).Error
}

func (u userRepository) Update(user *entities.User) error {
	return u.database.Model(&user).Where(entities.User{Username: user.Username}).Update(&user).Error
}

func (u userRepository) Get(username string) (*entities.User, error) {
	var user entities.User
	if err := u.database.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
