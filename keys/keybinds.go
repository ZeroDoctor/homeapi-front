package keys

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/views"
)

/**
* keys : code -> function
	u : 117 -> upload file/folder
	o : 111 -> opens file/folder
	c : 99  -> create file/folder
	d : 100 -> deletes file/folder
	g : 103 -> (gets) downloads file/folder
	m : 109 -> moves file/folder

	b : 98  -> pops state off of stack
	y : 121 -> accepts action
	n : 110 -> cancels action

	h : 104 -> moves cursor left
	j : 106 -> moves cursor down
	k : 107 -> moves cursor up
	l : 108 -> moves cursor right
*/

func setViewBindings(g *gocui.Gui, keys []KeyOp, viewName string) {
	for _, key := range keys {
		if err := g.SetKeybinding(viewName, key.keycode, gocui.ModNone, key.operation); err != nil {
			log.Panicln(err)
		}
	}
}

// SetGeneralBindings :
func SetGeneralBindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, views.Quit); err != nil { // instead of quit pop folder stack if empty then quit
		log.Panicln(err)
	}
	if err := g.SetKeybinding("tree", gocui.KeyTab, gocui.ModNone, views.NextView); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("screen", gocui.KeyTab, gocui.ModNone, views.NextView); err != nil {
		log.Panicln(err)
	}

	up := NewKey(rune(107), CursorUp)
	down := NewKey(rune(106), CursorDown)

	openTree := NewKey(rune(111), OpenTreeFile)

	openScreen := NewKey(rune(111), OpenScreenFile)
	downloadScreen := NewKey(rune(103), DownloadScreenFile)
	deleteScreen := NewKey(rune(100), DeleteScreenFile)
	uploadScreen := NewKey(rune(117), UploadScreenFile)

	treeKeys := []KeyOp{
		up,
		down,
		openTree,
	}
	screenKeys := []KeyOp{
		up,
		down,
		openScreen,
		downloadScreen,
		deleteScreen,
		uploadScreen,
	}

	setViewBindings(g, treeKeys, "tree")
	setViewBindings(g, screenKeys, "screen")
}
