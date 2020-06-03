package views

import (
	"fmt"
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
		fmt.Fprintln(v, "Press Enter to Start...")
		if _, err := SetCurrentViewOnTop(g, "status"); err != nil {
			return err
		}
	}

	return nil
}

// Logging :
func Logging(ftype string, msg interface{}, clear bool) Data {
	if ftype == "" {
		return NewData("", 0, clear, msg, nil)
	} else {
		return NewData(ftype, 0, clear, msg, nil)
	}	
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
				header = "\x1b[1;30;41m"+data.Type+"\x1b[0;37;40m : "
			case "Info": 
				header = "\x1b[1;30;42m"+data.Type+"\x1b[0;37;40m : "
			case "Warning":
				header = "\x1b[1;30;43m"+data.Type+"\x1b[0;37;40m : "
			}

			fmt.Fprint(statusView, header)
			fmt.Fprint(statusView, data.String)
		}
	}
}
