package ui

import (
	"fmt"

	c "github.com/tiagorlampert/CHAOS/src/color"
	"github.com/tiagorlampert/CHAOS/src/version"
)

func ShowMenu() {
	fmt.Println("")
	fmt.Println(c.YELLOW, "                                     CHAOS v"+version.GetVersion())
	fmt.Println(c.CYAN, "                               by tiagorlampert")
	fmt.Println("")
	fmt.Println(c.WHITE, " Please use `tab` to autocomplete commands,")
	fmt.Println(c.WHITE, " or type `exit` to quit this program.")
	fmt.Println("")
}
