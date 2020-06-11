package views

import (
	"sync"

	"github.com/jroimartin/gocui"
)

// SetupViewManager :
func SetupViewManager(g *gocui.Gui, wg *sync.WaitGroup) {
	wg.Add(3)
	go PrintTreeView(g, wg)
	go PrintScreenView(g, wg)
	go PrintStatusView(g, wg)
}
