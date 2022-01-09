package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func (client *UI) renderLatestView() fyne.CanvasObject {
	content := container.NewVBox()
	for _, c := range client.db.Channels {
		messages := client.db.ChanList[c.IdentHash]
		if len(messages) == 0 {
			continue
		}
		msg := messages[len(messages)-1]
		right := client.msgToCard(msg)
		content.Add(right)
	}
	navBar := client.NewNavbar("latest")
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, container.NewScroll(content))
}
