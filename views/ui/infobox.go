package ui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// SetInfoView :
func SetInfoView(g *gocui.Gui, screenView, emptyView *gocui.View, msg string, conX, x0, y1 int) error {
	newX := ((conX / 2) + (x0 / 2)) - 11
	newY := y1 / 2
	if v, err := g.SetView("empty", (newX - 1), newY-1, newX+20, newY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		screenView.Clear()
		v.Wrap = true
		v.Autoscroll = false
		fmt.Fprintln(v, msg)
		emptyView = v
		g.SetViewOnTop("empty")
	}

	return nil
}
