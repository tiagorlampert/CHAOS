package handler

type Handler interface {
	Handle() error
	Write(v string) error
}
