package handler

type Handler interface {
	Handle() error
	WriteCommandResponse(input string)
	Write(v string) error
}
