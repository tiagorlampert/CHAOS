package environment

import "os"

type Configuration struct {
	Server   Server
	Database Database
}

type Server struct {
	Port string `validate:"required"`
}

type Database struct {
	Name string `validate:"required"`
}

func LoadEnv() *Configuration {
	return &Configuration{
		Server: Server{
			Port: os.Getenv(`PORT`),
		},
		Database: Database{
			Name: os.Getenv(`DATABASE_NAME`),
		},
	}
}
