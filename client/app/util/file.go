package util

import (
	"io/ioutil"
	"os"
)

func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func WriteFile(path string, data []byte) error {
	if err := ioutil.WriteFile(path, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func CopyFile(source, destination string) error {
	data, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(destination, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}
