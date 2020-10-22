package models

type Message struct {
	Command string
	Data    []byte
	Error   Error
}

type Error struct {
	HasError bool
	Message  string
}
