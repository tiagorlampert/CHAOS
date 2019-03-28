package cmd

import (
	"github.com/tiagorlampert/CHAOS/src/ui"
	"github.com/tiagorlampert/CHAOS/src/util"
)

func Start() {
	util.DetectOS()
	util.ClearScreen()
	ui.StartMenu()
}
