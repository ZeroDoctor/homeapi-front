package views

import (
	"fmt"
	"log"
	"sync"

	"github.com/jroimartin/gocui"
)

var statusView *gocui.View

// SetStatusView :
func SetStatusView(g *gocui.Gui, maxX int, maxY int) error {

	if v, err := g.SetView("status", 0, maxY-(maxY/15)-1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "status"
		v.Wrap = false

		statusView = v
	}

	return nil
}

// PrintStatusView :
func PrintStatusView(g *gocui.Gui, wg *sync.WaitGroup) {
	defer wg.Done()
	for data := range InStatusChan {
		if statusView != nil {
			if data.Boolean {
				statusView.Clear()
			}

			var header interface{}
			header = ""
			switch data.Type {
			case "Error":
				header = "\x1b[1;30;41mError\x1b[0;37;40m : "
			case "Info":
				header = "\x1b[1;30;42mInfo\x1b[0;37;40m : "
			case "Warning":
				header = "\x1b[1;30;43mWarning\x1b[0;37;40m : "
			case "delete": // move would also be in this case
				header = "\x1b[1;30;43m" + data.Type + "\x1b[0;37;40m : "
				questionStatusView(g, data)
			}

			fmt.Fprint(statusView, header)
			fmt.Fprint(statusView, data.String)
		}
	}
}

// Logging :
func Logging(ftype string, msg interface{}, clear bool) Data {
	if ftype == "" {
		return NewData("", 0, clear, msg, nil)
	}

	return NewData(ftype, 0, clear, msg, nil)
}

func questionStatusView(g *gocui.Gui, data Data) {

	yErr := g.SetKeybinding("status", rune(121), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { // 'y' key
		SetCurrentViewOnTop(g, "screen")
		InScreenChan <- data
		return nil
	})

	if yErr != nil {
		log.Panicln(yErr)
	}

	nErr := g.SetKeybinding("status", rune(110), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { // 'n' key
		SetCurrentViewOnTop(g, "screen")
		InScreenChan <- NewData("cancel", 1, false, "", nil)
		return nil
	})

	if nErr != nil {
		log.Panicln(nErr)
	}
}
