package ui

import "fmt"

func ShowMenu(version string, port string) {
	fmt.Println(" ┌───────────────────────────────────────────────────┐ ")
	fmt.Printf(" │                    CHAOS %s                     │ \n", version)
	fmt.Printf(" │             http://127.0.0.1:%s                 │ \n", port)
	fmt.Println(" │                by tiagorlampert                   │ ")
	fmt.Println(" └───────────────────────────────────────────────────┘")
	fmt.Println("")
}
