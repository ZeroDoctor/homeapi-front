package keys

import (
	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/views"
)

// CursorUp :
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

// CursorDown :
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// CursorLeft : reserved for files
func CursorLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {

	}
	return nil
}

// CursorRight : reserved for files
func CursorRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {

	}
	return nil
}

// InitTree :
func InitTree(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, err := g.View("tree")
		if err != nil {
			return err
		}

		if _, err := views.SetCurrentViewOnTop(g, "tree"); err != nil {
			return err
		}

		views.InTreeChan <- views.NewData("init", 0, false, "", nil)
	}
	return nil
}

// OpenTreeFile :
func OpenTreeFile(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		views.InTreeChan <- views.NewData("open", cy+oy, false, "", nil)
	}
	return nil
}

// OpenScreenFile :
func OpenScreenFile(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()
		views.InScreenChan <- views.NewData("view", cy+oy-1, false, "", nil)
	}

	return nil
}

func DownloadScreenFile(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()
		_, oy := v.Origin()

		views.InScreenChan <- views.NewData("download", cy+oy-2, false, "", nil)
	}

	return nil
}

func DeleteScreenFile(g *gocui.Gui, v *gocui.View) error {

	return nil
}
