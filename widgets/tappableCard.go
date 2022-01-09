package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/data"
)

type TappableCard struct {
	widget.Card
	msg             *data.Message
	ChanID          string
	SelectedMessage chan int
	SelectedChannel chan string
}

func NewTappableCard(title, subtitle string, content fyne.CanvasObject) *TappableCard {
	card := &TappableCard{}
	card.ExtendBaseWidget(card)
	card.SetTitle(title)
	card.SetSubTitle(subtitle)
	card.SetContent(content)
	return card
}

func (t *TappableCard) Tapped(_ *fyne.PointEvent) {
	if t.SelectedMessage != nil {
		t.SelectedMessage <- t.msg.ID
		close(t.SelectedMessage)
	}
	if t.SelectedChannel != nil {
		t.SelectedChannel <- t.ChanID
		close(t.SelectedChannel)
	}
}

func (t *TappableCard) TappedSecondary(_ *fyne.PointEvent) {
}
