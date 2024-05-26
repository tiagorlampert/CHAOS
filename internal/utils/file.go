package utils

import "os"

func WriteFile(filepath string, s []byte) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(s)
	return err
}
