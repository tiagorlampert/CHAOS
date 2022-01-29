package system

import "os"

// CreateDirs create a directory based on a given paths
func CreateDirs(paths ...string) error {
	for _, path := range paths {
		if err := createDirectory(path); err != nil {
			return err
		}
	}
	return nil
}

func createDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}
