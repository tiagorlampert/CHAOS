package ui

import (
	"fmt"
	c "github.com/tiagorlampert/CHAOS/pkg/color"
)

func ShowMenuHeader(version string) {
	fmt.Println("")
	fmt.Println(c.Yellow, "                                                 CHAOS "+version)
	fmt.Println(c.Cyan, "                                             by tiagorlampert")
	fmt.Println("")
	fmt.Println(c.White, " Please use `tab` to autocomplete commands,")
	fmt.Println(c.White, " or type `exit` to quit this program.")
	fmt.Println("")
}
