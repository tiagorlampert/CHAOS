package device

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"gorm.io/gorm"
	"time"
)

type deviceRepository struct {
	dbClient *gorm.DB
}

func NewRepository(dbClient *gorm.DB) Repository {
	return &deviceRepository{dbClient: dbClient}
}

func (r deviceRepository) Insert(input entities.Device) error {
	result := r.dbClient.Create(&input)
	if result.Error != nil {
		return repositories.HandleError(result.Error)
	}
	if result.RowsAffected <= 0 {
		return errors.New("error saving device")
	}
	return nil
}

func (r deviceRepository) Update(device entities.Device) error {
	return r.dbClient.Model(&device).Where(
		entities.Device{MacAddress: device.MacAddress}).Updates(&device).Error
}

func (r deviceRepository) FindByMacAddress(address string) (*entities.Device, error) {
	var device entities.Device
	if err := r.dbClient.Where(entities.Device{MacAddress: address}).First(&device).Error; err != nil {
		return nil, repositories.HandleError(err)
	}
	return &device, nil
}

func (r deviceRepository) FindAll(updatedAt time.Time) ([]entities.Device, error) {
	var devices []entities.Device
	if err := r.dbClient.Where(
		"updated_at > ?", updatedAt).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
