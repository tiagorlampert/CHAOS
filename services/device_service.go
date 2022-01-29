package services

import (
	"errors"
	"github.com/tiagorlampert/CHAOS/entities"
	"github.com/tiagorlampert/CHAOS/internal/utilities"
	"github.com/tiagorlampert/CHAOS/repositories"
	"time"
)

type deviceService struct {
	repository repositories.Device
}

func NewDevice(repository repositories.Device) Device {
	return &deviceService{repository: repository}
}

func (d deviceService) Insert(input entities.Device) error {
	_, err := d.repository.GetByMacAddress(input.MacAddress)
	if errors.Is(err, repositories.ErrNotFound) {
		return d.repository.Insert(input)
	}
	return d.repository.Update(input)
}

func (d deviceService) FindAll() ([]entities.Device, error) {
	devices, err := d.repository.FindAll(time.Now().Add(time.Minute * time.Duration(-3)))
	if err != nil {
		return nil, err
	}
	for index, device := range devices {
		devices[index].MacAddressBase64 = utilities.EncodeBase64(device.MacAddress)
	}
	return devices, nil
}
