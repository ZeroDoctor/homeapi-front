package keys

import "github.com/jroimartin/gocui"

// KeyOp :
type KeyOp struct {
	keycode   rune
	operation func(*gocui.Gui, *gocui.View) error
}

// NewKey :
func NewKey(keycode rune, operation func(*gocui.Gui, *gocui.View) error) KeyOp {
	return KeyOp{
		keycode:   keycode,
		operation: operation,
	}
}
