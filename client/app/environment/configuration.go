package environment

import "time"

type Configuration struct {
	Connection Connection
	Server     Server
}

type Connection struct {
	Token           string
	ContextDeadline time.Duration
}

type Server struct {
	Address  string
	HttpPort string
	Url      string
}
