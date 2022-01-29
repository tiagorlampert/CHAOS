package system

import "os"

// CreateDir create a directory based on a given path
func CreateDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

// CreateDirs create a directory based on a given paths
func CreateDirs(paths ...string) error {
	for _, path := range paths {
		if err := CreateDir(path); err != nil {
			return err
		}
	}
	return nil
}
