package main

import "fmt"

type SamplePlugin struct{}

func (p *SamplePlugin) Initialize() error {
    fmt.Println("Sample plugin initialized")
    return nil
}

func (p *SamplePlugin) Execute() error {
    fmt.Println("Sample plugin executed")
    return nil
}

func (p *SamplePlugin) Cleanup() error {
    fmt.Println("Sample plugin cleaned up")
    return nil
}

// This is important! Go plugins must export a variable named "Plugin"
var Plugin SamplePlugin