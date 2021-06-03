package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (client *GUI) renderLatestView() fyne.CanvasObject {
	content := container.New(layout.NewFormLayout())
	for _, c := range client.db.Channels {
		messages := client.db.chanList[c.IdentHash]
		if len(messages) == 0 {
			continue
		}
		msg := messages[len(messages)-1]
		chanCard := widget.NewCard("", client.db.nameFromChanIdentHash(msg.TargetChannel), nil)
		img := client.avatarCache[c.IdentHash]
		img.SetMinSize(fyne.NewSize(64, 64))
		img.FillMode = canvas.ImageFillStretch
		chanCard.SetImage(img)
		content.Add(chanCard)

		card := client.msgToCard(msg)
		content.Add(card)
	}
	navBar := client.renderNavBar("latest")
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewVScroll(content))
}
