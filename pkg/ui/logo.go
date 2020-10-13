package ui

import (
	"fmt"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
)

func ShowLogo() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println(c.Red, "▄████████    ▄█    █▄       ▄████████  ▄██████▄     ▄████████  ")
	fmt.Println(c.Red, "███    ███   ███    ███     ███    ███ ███    ███   ███    ███ ")
	fmt.Println(c.Red, "███    █▀    ███    ███     ███    ███ ███    ███   ███    █▀  ")
	fmt.Println(c.Red, "███         ▄███▄▄▄▄███▄▄   ███    ███ ███    ███   ███        ")
	fmt.Println(c.Red, "███        ▀▀███▀▀▀▀███▀  ▀███████████ ███    ███ ▀███████████ ")
	fmt.Println(c.Red, "███    █▄    ███    ███     ███    ███ ███    ███          ███ ")
	fmt.Println(c.Red, "███    ███   ███    ███     ███    ███ ███    ███    ▄█    ███ ")
	fmt.Println(c.Red, "████████▀    ███    █▀      ███    █▀   ▀██████▀   ▄████████▀  ")
	fmt.Println("")
}
