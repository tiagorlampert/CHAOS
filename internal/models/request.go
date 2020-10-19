package models

type Request struct {
	Runnable bool
	Command  string
	Data     []byte
}
