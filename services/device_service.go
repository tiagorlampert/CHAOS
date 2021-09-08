package services

import (
	"github.com/tiagorlampert/CHAOS/entities"
	repo "github.com/tiagorlampert/CHAOS/repositories"
	"github.com/tiagorlampert/CHAOS/shared/utils"
	"time"
)

type deviceService struct {
	repository repo.Device
}

func NewDevice(repository repo.Device) Device {
	return &deviceService{
		repository: repository,
	}
}

func (d deviceService) Insert(input entities.Device) error {
	_, err := d.repository.Get(input.MacAddress)
	if err != nil {
		if err == repo.ErrNotFound {
			if err := d.repository.Insert(input); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if err := d.repository.Update(input); err != nil {
		return err
	}
	return nil
}

func (d deviceService) GetAllAvailable() ([]entities.Device, error) {
	devices, err := d.repository.List(utils.GetTimeWith(time.Minute, -3))
	if err != nil {
		return nil, err
	}

	for pos, device := range devices {
		devices[pos].MacAddressBase64 = utils.EncodeBase64(device.MacAddress)
	}
	return devices, nil
}
