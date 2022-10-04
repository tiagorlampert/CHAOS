package repositories

import (
	"errors"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

func HandleError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return ErrNotFound
	default:
		return err
	}
}
