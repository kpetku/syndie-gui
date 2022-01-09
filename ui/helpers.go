package ui

import (
	"log"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/data"
	"github.com/kpetku/syndie-gui/widgets"
)

func newCenteredContainer(l *fyne.Container) *fyne.Container {
	if !fyne.CurrentDevice().IsMobile() {
		// Center the layout on desktop
		return container.New(layout.NewGridLayoutWithColumns(1), l)
	}
	return l
}

func (client UI) msgToCard(msg data.Message) *widget.Card {
	date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
	hbox := container.New(layout.NewFormLayout())
	vbox := container.NewVBox()
	authorAvatar := client.avatarCache[msg.Author]
	if authorAvatar == nil {
		authorAvatar = &canvas.Image{}
	}
	authorAvatar.SetMinSize(fyne.NewSize(32, 32))
	authorAvatar.FillMode = canvas.ImageFillOriginal
	hbox.Add(authorAvatar)
	hbox.Add(widget.NewLabel("by " + client.db.NameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02")))
	vbox.Add(hbox)
	if len(msg.Raw.Page) > 0 {
		for num, p := range msg.Raw.Page[:1] {
			if num >= 0 {
				vbox.Add(widgets.NewLabel("Page: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(msg.Raw.Page)-1)))
				vbox.Add(widget.NewSeparator())
				mu, err := widgets.NewMarkup(msg, p.Data)
				if err != nil {
					log.Printf("NewMarkup err: %s", err.Error())
				}
				vbox.Add(mu)
			}
		}
	}
	return widget.NewCard(msg.Subject, "posted to "+client.db.NameFromChanIdentHash(msg.TargetChannel)+" "+shortIdent(msg.TargetChannel), vbox)
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}

func (client UI) pagination() int {
	i, err := strconv.Atoi(client.app.Preferences().String("pagination"))
	if err != nil {
		i = 25
	}
	return i
}
