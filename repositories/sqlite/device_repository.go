package sqlite

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"time"
)

type deviceRepository struct {
	database *gorm.DB
}

func NewDeviceRepository(database *gorm.DB) repositories.Device {
	return &deviceRepository{database: database}
}

func (r deviceRepository) Insert(input entities.Device) error {
	rowsAffected := r.database.Create(&input).RowsAffected
	if rowsAffected <= 0 {
		return errors.New("error saving device")
	}
	return nil
}

func (r deviceRepository) Update(device entities.Device) error {
	return r.database.Model(&device).Where(entities.Device{MacAddress: device.MacAddress}).Update(&device).Error
}

func (r deviceRepository) Get(macAddress string) (*entities.Device, error) {
	var device entities.Device
	err := r.database.Where(entities.Device{MacAddress: macAddress}).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}
	return &device, nil
}

func (r deviceRepository) List(dateTime time.Time) ([]entities.Device, error) {
	var devices []entities.Device
	if err := r.database.Where("updated_at > ?", dateTime.String()).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
