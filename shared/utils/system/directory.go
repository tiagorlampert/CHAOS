package system

import "os"

// CreateDirectory create a directory based on a given path
func CreateDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
