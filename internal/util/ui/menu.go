package ui

import (
	"github.com/tiagorlampert/CHAOS/pkg/system"
	"github.com/tiagorlampert/CHAOS/pkg/ui"
)

func ShowMenu(version string) {
	system.ClearScreen()
	ui.ShowLogo()
	ui.ShowMenuHeader(version)
}
