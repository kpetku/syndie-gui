package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/data"
)

type tappableCard struct {
	widget.Card
	msg             *data.Message
	chanID          string
	selectedMessage chan int
	selectedChannel chan string
}

func newTappableCard(title, subtitle string, content fyne.CanvasObject) *tappableCard {
	card := &tappableCard{}
	card.ExtendBaseWidget(card)
	card.SetTitle(title)
	card.SetSubTitle(subtitle)
	card.SetContent(content)
	return card
}

func (t *tappableCard) Tapped(_ *fyne.PointEvent) {
	if t.selectedMessage != nil {
		t.selectedMessage <- t.msg.ID
		close(t.selectedMessage)
	}
	if t.selectedChannel != nil {
		t.selectedChannel <- t.chanID
		close(t.selectedChannel)
	}
}

func (t *tappableCard) TappedSecondary(_ *fyne.PointEvent) {
}

func newLabel(s string) *widget.Label {
	w := widget.NewLabel(s)
	w.Wrapping = fyne.TextWrapBreak
	return w
}
