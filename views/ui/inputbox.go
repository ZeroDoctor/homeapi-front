package ui

import (
	"fmt"
	"regexp"

	"github.com/jroimartin/gocui"
)

// InputBox :
var InputBox = &Input{}
var screenChan chan interface{}

// Input :
type Input struct {
	str []rune
}

// SetInputBox :
func SetInputBox(g *gocui.Gui, dialogView *gocui.View, title string, conX, x0, y1 int) error {
	newX := ((conX / 2) + (x0 / 2)) - 11
	newY := y1 / 2
	if v, err := g.SetView("dialog", (newX - 1), newY-1, newX+20, newY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		v.Autoscroll = true
		v.Editable = true
		v.Editor = InputBox
		v.Title = title
		fmt.Fprintln(v, "C:/")
		dialogView = v
	}

	return nil
}

// TabComplete : this is going somewhere else later
func TabComplete(text string, list []string) ([]string, error) {
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

// Edit :
func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		i.str = append(i.str, ch)
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		i.str = append(i.str, rune(' '))
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}
