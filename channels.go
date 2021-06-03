package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/data"
)

func (client *GUI) renderThreadList(needle int) *fyne.Container {
	client.threadPane = container.NewVBox()
	for num, msg := range client.db.ChanList[client.selectedChannel] {
		if num <= client.channelNeedle || num < client.pagination {
			currentMessage := msg
			client.threadPane.Add(client.msgToCard(currentMessage))
		}
		if num == client.channelNeedle {
			if num <= client.pagination {
				continue
			}
			client.threadPane.Add(widget.NewButton("Show more messages", func() {
				client.channelNeedle = needle + client.pagination
				client.renderThreadList(needle + client.pagination)
				client.repaintMainWindow()
			}))
		}
	}
	if needle == client.channelNeedle {
		if len(client.db.ChanList[client.selectedChannel])-1 <= client.pagination {
			return client.threadPane
		}
		if needle <= client.pagination {
			client.threadPane.Add(widget.NewButton("Show more messages", func() {
				if client.channelNeedle == needle {
					client.channelNeedle = client.pagination + (client.pagination - 1)
				} else {
					client.channelNeedle = needle + client.pagination
				}
				client.renderThreadList(needle + client.pagination)
				client.repaintMainWindow()
			}))
		}
	}
	return client.threadPane
}

func (client *GUI) renderThreadListWithMenu(needle int) *fyne.Container {
	backButton := widget.NewButton("Back", func() {
		client.selectedChannel = ""
		client.repaintMainWindow()
	})
	navBar := container.NewGridWithColumns(3)
	navBar.Add(backButton)
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewVScroll(client.renderThreadList(needle)))
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}

func (client GUI) msgToCard(msg data.Message) *widget.Card {
	date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
	text := "by " + client.db.NameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02")
	vbox := container.NewVBox()

	hbox := container.NewHBox()
	icon := client.avatarCache[msg.Author]
	if icon != nil {
		hbox.Add(icon)
	}
	hbox.Add(widget.NewLabel(text))

	vbox.Add(hbox)
	if len(msg.Raw.Page) > 0 {
		for num, p := range msg.Raw.Page[:1] {
			if num >= 0 {
				vbox.Add(newLabel("Page: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(msg.Raw.Page)-1)))
				vbox.Add(widget.NewSeparator())
				vbox.Add(newLabel(p.Data))
			}
		}
	}
	return widget.NewCard(msg.Subject, "", vbox)
}
