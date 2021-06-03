package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewLabel(s string) *widget.Label {
	w := widget.NewLabel(s)
	w.Wrapping = fyne.TextWrapBreak
	return w
}
