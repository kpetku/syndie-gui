package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewLabel(s string) *widget.RichText {
	w := widget.NewRichTextWithText(s)
	w.Wrapping = fyne.TextWrapBreak
	return w
}
