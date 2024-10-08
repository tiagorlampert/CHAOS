package entities

type Plugin interface {
    Initialize() error
    Execute() error
    Cleanup() error
}