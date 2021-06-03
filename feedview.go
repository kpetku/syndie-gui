package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const noDescription = "No description"
const noName = "Anonymous"

func (client *GUI) renderFeedView() fyne.CanvasObject {
	content := container.New(layout.NewFormLayout())
	for _, c := range client.db.Channels {
		var desc, name string
		if c.Description == "" {
			desc = noDescription
		} else {
			desc = c.Description
		}
		if c.Name == "" {
			name = noName
		} else {
			name = c.Name
		}
		img := client.avatarCache[c.IdentHash]
		img.SetMinSize(fyne.NewSize(64, 64))
		img.FillMode = canvas.ImageFillContain
		content.Add(img)

		card := newTappableCard(name, c.IdentHash, widget.NewLabel(desc))
		card.chanID = c.IdentHash
		card.selectedChannel = make(chan string)
		go func() {
			for click := range card.selectedChannel {
				c := click
				client.selectedChannel = c
				break
			}
			go client.repaintMainWindow()
		}()
		content.Add(card)
	}
	navBar := client.renderNavBar()
	return newCenteredContainer(container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewVScroll(content)))
}
