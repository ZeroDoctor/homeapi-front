package main

import (
	"log"
	"regexp"
	"sync"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/keys"
	"github.com/zerodoctor/homeapi-front/views"
)

// ?: this is going somewhere later
func tabComplete(text string, list []string) ([]string, error) {
	var possibles []string

	for _, word := range list {
		match, err := regexp.MatchString(`(?i)^(`+text+`)`, word)
		if err != nil {
			return nil, err
		}
		if match {
			possibles = append(possibles, word)
		}
	}

	return possibles, nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	err := views.SetTreeView(g, maxX, maxY)
	if err != nil {
		return err
	}
	err = views.SetScreenView(g, maxX, maxY)
	if err != nil {
		return err
	}
	err = views.SetStatusView(g, maxX, maxY)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Mouse = false
	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorCyan

	init := []string{"tree", "screen"}
	views.NewViews(init)
	var wg sync.WaitGroup
	views.SetupViewManager(g, &wg)

	g.SetManagerFunc(layout)
	keys.SetGeneralBindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
