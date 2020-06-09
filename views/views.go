package views

import (
	"github.com/jroimartin/gocui"
)

var vl *ViewList

// ViewList :
type ViewList struct {
	viewList    []string
	currentView int
}

// NewViews :
func NewViews(list []string) {
	vl = &ViewList{list, 0}
}

// SetCurrentViewOnTop :
func SetCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

// NextView :
func NextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (vl.currentView + 1) % (len(vl.viewList))
	name := vl.viewList[nextIndex]

	if _, err := SetCurrentViewOnTop(g, name); err != nil {
		return err
	}

	vl.currentView = nextIndex
	return nil
}

// Quit :
func Quit(g *gocui.Gui, v *gocui.View) error {
	close(InTreeChan)
	close(InScreenChan)
	close(InStatusChan)

	return gocui.ErrQuit
}
