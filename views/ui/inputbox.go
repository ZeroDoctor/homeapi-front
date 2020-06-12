package ui

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/jroimartin/gocui"
	"github.com/zerodoctor/homeapi-front/channel"
)

// InputBox :
var InputBox = &Input{}
var inputView *gocui.View

// Input :
type Input struct {
	Str []rune
}

// SetInputBox :
func SetInputBox(g *gocui.Gui, title string, conX, x0, x1, y1 int) (*gocui.View, error) {
	width := ((x1 - x0) / 2) - 1
	newX := ((conX / 2) + (x0 / 2)) - (width / 2)
	newY := y1 / 2
	if v, err := g.SetView("dialog", (newX - 1), newY-1, newX+width, newY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		v.Wrap = false
		v.Autoscroll = true
		v.Editable = true
		InputBox.Str = []rune("C:/")
		v.Editor = InputBox
		v.Title = title
		inputView = v

		for _, c := range InputBox.Str {
			v.EditWrite(c)
		}
	}

	return inputView, nil
}

// Edit :
func (i *Input) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		i.Str = append(i.Str, ch)
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		i.Str = append(i.Str, rune(' '))
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyTab:
		autoComplete(v)
	case key == gocui.KeyEnd:
		channel.InScreenChan <- channel.NewData("cancel", 0, false, "", nil)
	}
}

func autoComplete(v *gocui.View) {
	if InputBox.Str == nil {
		return
	}

	//path := strings.Split(string(InputBox.Str), "/")

}

// getFilesFolder :
func getFilesFolder(path string) ([]string, error) {
	var result []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		result = append(result, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// tabComplete : this is going somewhere else later
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
