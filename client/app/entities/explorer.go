package entities

import "time"

type FileExplorer struct {
	Path        string   `json:"path"`
	Files       []File   `json:"files"`
	Directories []string `json:"directories"`
}

type File struct {
	Filename string    `json:"filename"`
	ModTime  time.Time `json:"mod_time"`
}
