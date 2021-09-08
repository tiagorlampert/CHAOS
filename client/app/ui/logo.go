package ui

import "fmt"

func ShowLogo(version string) {
	if version == "dev" {
		version = fmt.Sprint(version, " ")
	}
	fmt.Println(" ┌───────────────────────────────────────────────────┐ ")
	fmt.Printf(" │                   CHAOS %s                      │ \n", version)
	fmt.Println(" │                by tiagorlampert                   │ ")
	fmt.Println(" └───────────────────────────────────────────────────┘")
}
