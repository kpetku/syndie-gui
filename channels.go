package main

import (
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func (client *GUI) renderChannelList() fyne.CanvasObject {
	right := widget.NewVBox()
	for _, channel := range client.db.Channels {
		row := widget.NewHBox()
		numOfMessages := strconv.Itoa(len(client.db.chanList[channel.IdentHash]))
		test := canvas.NewImageFromImage(client.db.getAvatar(channel.IdentHash))
		test.SetMinSize(fyne.NewSize(32, 32))
		test.FillMode = canvas.ImageFillContain
		row.Append(test)
		row.Append(widget.NewLabel(channel.Name + " " + shortIdent(channel.IdentHash) + " " + numOfMessages + "/" + numOfMessages))
		right.Append(row)
	}
	return right
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}
