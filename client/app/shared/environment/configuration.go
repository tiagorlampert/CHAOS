package environment

import "time"

type Configuration struct {
	Connection Connection
	Server     Server
}

type Connection struct {
	Token             string
	ContextDeadline   time.Duration
	ContentTypeHeader string
	ContentTypeJSON   string
	CookieHeader      string
}

type Server struct {
	Address string
	Port    string
	URL     string
	Endpoint
}

type Endpoint struct {
	Health   string
	Device   string
	Command  string
	Upload   string
	Download string
}
