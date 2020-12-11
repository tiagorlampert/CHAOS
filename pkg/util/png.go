package util

import (
	"io/ioutil"
	"os"
)

// WriteFile create a new file
func WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, os.ModePerm)
}
