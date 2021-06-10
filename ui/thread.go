package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (client *UI) renderThreadList(needle int) *fyne.Container {
	client.threadPane = container.NewVBox()
	for num, msg := range client.db.ChanList[client.selectedChannel] {
		if num <= client.channelNeedle || num < client.pagination() {
			currentMessage := msg
			client.threadPane.Add(client.msgToCard(currentMessage))
		}
		if num == client.channelNeedle {
			if num <= client.pagination() {
				continue
			}
			client.threadPane.Add(widget.NewButton("Show more messages", func() {
				client.channelNeedle = needle + client.pagination()
				client.renderThreadList(needle + client.pagination())
				client.repaintMainWindow()
			}))
		}
	}
	if needle == client.channelNeedle {
		if len(client.db.ChanList[client.selectedChannel])-1 <= client.pagination() {
			return client.threadPane
		}
		if needle <= client.pagination() {
			client.threadPane.Add(widget.NewButton("Show more messages", func() {
				if client.channelNeedle == needle {
					client.channelNeedle = client.pagination() + (client.pagination() - 1)
				} else {
					client.channelNeedle = needle + client.pagination()
				}
				client.renderThreadList(needle + client.pagination())
				client.repaintMainWindow()
			}))
		}
	}
	return client.threadPane
}

func (client *UI) renderThreadListWithMenu(needle int) *fyne.Container {
	backButton := widget.NewButton("Back", func() {
		client.selectedChannel = ""
		client.repaintMainWindow()
	})
	navBar := container.NewGridWithColumns(3)
	navBar.Add(backButton)
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewVScroll(client.renderThreadList(needle)))
}
