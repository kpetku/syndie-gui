package main

import (
	"github.com/kpetku/syndie-core/data"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type tappableLabel struct {
	widget.Label
	msg             *data.Message
	chanID          string
	selectedMessage chan int
	selectedChannel chan string
}

func newTappableLabel(text string) *tappableLabel {
	label := &tappableLabel{}
	label.ExtendBaseWidget(label)
	label.SetText(text)
	label.Wrapping = fyne.TextWrapBreak
	return label
}

func (t *tappableLabel) Tapped(_ *fyne.PointEvent) {
	if t.selectedMessage != nil {
		t.selectedMessage <- t.msg.ID
		close(t.selectedMessage)
	}
	if t.selectedChannel != nil {
		t.selectedChannel <- t.chanID
		close(t.selectedChannel)
	}
}

func (t *tappableLabel) TappedSecondary(_ *fyne.PointEvent) {
}

func newLabel(s string) *widget.Label {
	w := widget.NewLabel(s)
	w.Wrapping = fyne.TextWrapBreak
	return w
}
