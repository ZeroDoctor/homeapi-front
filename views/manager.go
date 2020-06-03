package views

import (
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/model"
)

// Data : this give the ability to send mutliple data types not sure if this is right
type Data struct {
	Type    string
	Integer int
	Boolean bool
	String  interface{}
	File    []model.FileFolder
}

// NewData :
func NewData(dtype string, integer int, boolean bool, myString interface{}, file []model.FileFolder) Data {
	return Data{
		Type:    dtype,
		Integer: integer,
		Boolean: boolean,
		String:  myString,
		File:    file,
	}
}

// Is this a lot of channels?
var (
	// InTreeChan : could be int instead of Data, but maybe in the future I could need a string and int
	InTreeChan = make(chan Data, 4)
	// OutTreeChan :
	OutTreeChan = make(chan Data)

	// InScreenChan :
	InScreenChan = make(chan Data, 4)
	// OutScreenChan :
	OutScreenChan = make(chan Data)

	// InStatusChan :
	InStatusChan = make(chan Data, 2)
	// OutStatusChan :
	OutStatusChan = make(chan Data)
)

// SetupViewManager :
func SetupViewManager(g *gocui.Gui, wg *sync.WaitGroup) {
	wg.Add(3)
	go PrintTreeView(g, wg)
	go PrintScreenView(g, wg)
	go PrintStatusView(g, wg)
}
