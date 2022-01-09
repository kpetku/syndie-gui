package ui

import (
	"sort"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-gui/widgets"
)

const noDescription = "No description"
const noName = "Anonymous"

func (client *UI) renderFeedView() fyne.CanvasObject {
	data := []*widgets.TappableCard{}
	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return widgets.NewTappableCard("template", "template", widget.NewLabel("template"))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widgets.TappableCard).ChanID = data[i].ChanID
			o.(*widgets.TappableCard).SetTitle(data[i].Title)
			o.(*widgets.TappableCard).SetSubTitle(data[i].Subtitle)
			o.(*widgets.TappableCard).SetContent(data[i].Content)
			o.(*widgets.TappableCard).SelectedChannel = data[i].SelectedChannel
			o.(*widgets.TappableCard).SelectedMessage = data[i].SelectedMessage
		})

	for _, c := range client.db.Channels {
		messages := client.db.ChanList[c.IdentHash]
		// We don't use timestamps, therefore we must sort messages from largest to smallest MessageID because that is the closest strategy for ordering
		sort.Slice(messages, func(i, j int) bool { return messages[i].ID > messages[j].ID })
		if client.app.Preferences().Bool("hideEmptyFeeds") && len(messages) == 0 {
			continue
		}
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
		hbox := container.New(layout.NewFormLayout())
		img := client.avatarCache[c.IdentHash]
		img.SetMinSize(fyne.NewSize(32, 32))
		img.FillMode = canvas.ImageFillContain
		hbox.Add(img)

		card := widgets.NewTappableCard(name, shortIdent(c.IdentHash), hbox)
		hbox.Add(widgets.NewLabel(desc))
		card.ChanID = c.IdentHash
		card.SelectedChannel = make(chan string)
		go func() {
			for click := range card.SelectedChannel {
				c := click
				client.selectedChannel = c
				break
			}
			go client.repaintMainWindow()
		}()
		data = append(data, card)
	}
	navBar := client.NewNavbar("feed")
	return container.New(layout.NewBorderLayout(navBar, nil, nil, nil), navBar, list)
}
