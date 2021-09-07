package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (client *UI) renderLatestView() fyne.CanvasObject {
	content := container.New(layout.NewFormLayout())
	for _, c := range client.db.Channels {
		messages := client.db.ChanList[c.IdentHash]
		if len(messages) == 0 {
			continue
		}
		msg := messages[len(messages)-1]
		left := container.NewVBox()
		left.Add(widget.NewLabel(client.db.NameFromChanIdentHash(msg.TargetChannel)))
		img := client.avatarCache[c.IdentHash]
		img.SetMinSize(fyne.NewSize(64, 64))
		img.FillMode = canvas.ImageFillContain
		left.Add(img)
		content.Add(left)
		right := client.msgToCard(msg)
		content.Add(right)
	}
	navBar := client.NewNavbar("latest")
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewScroll(content))
}
