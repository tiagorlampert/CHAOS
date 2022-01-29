package sqlite

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"time"
)

type deviceSqliteRepository struct {
	dbClient *gorm.DB
}

func NewDeviceRepository(dbClient *gorm.DB) repositories.Device {
	return &deviceSqliteRepository{dbClient: dbClient}
}

func (r deviceSqliteRepository) Insert(input entities.Device) error {
	result := r.dbClient.Create(&input)
	if result.Error != nil {
		return handleError(result.Error)
	}
	if result.RowsAffected <= 0 {
		return errors.New("error saving device")
	}
	return nil
}

func (r deviceSqliteRepository) Update(device entities.Device) error {
	return r.dbClient.Model(&device).Where(
		entities.Device{MacAddress: device.MacAddress}).Update(&device).Error
}

func (r deviceSqliteRepository) GetByMacAddress(address string) (*entities.Device, error) {
	var device entities.Device
	if err := r.dbClient.Where(entities.Device{MacAddress: address}).First(&device).Error; err != nil {
		return nil, handleError(err)
	}
	return &device, nil
}

func (r deviceSqliteRepository) FindAll(updatedAt time.Time) ([]entities.Device, error) {
	var devices []entities.Device
	if err := r.dbClient.Where(
		"updated_at > ?", updatedAt.String()).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func handleError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return repositories.ErrNotFound
	default:
		return err
	}
}
