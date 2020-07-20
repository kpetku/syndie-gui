package main

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

func (client *GUI) renderChannelList() fyne.CanvasObject {
	var hideEmptyChannels = true

	client.channelPane = widget.NewVBox()
	for _, channel := range client.db.Channels {
		if hideEmptyChannels {
			if len(client.db.chanList[channel.IdentHash]) == 0 {
				continue
			}
		}
		rw := new(ChannelWidget)
		rw.Name = channel.Name
		rw.Text = widget.NewLabel(channel.IdentHash)
		rw.Icon = canvas.NewImageFromImage(client.db.getAvatar(channel.IdentHash))

		rw.SelectedChannel = make(chan string)
		go func() {
			for contact := range rw.SelectedChannel {
				client.selectedChannel = contact
				close(rw.SelectedChannel)
				client.repaint()
			}
		}()
		client.channelPane.Append(rw)
	}
	return client.channelPane
}

func (client *GUI) renderThreadList() fyne.CanvasObject {
	client.threadPane = widget.NewVBox()
	if client.selectedChannel == "" {
		client.threadPane.Append(widget.NewLabel("Select a channel from the menu to the left to get started"))
	} else {
		for _, msg := range client.db.chanList[client.selectedChannel] {
			client.threadPane.Append(widget.NewLabel(msg.Subject))
			date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
			client.threadPane.Append(widget.NewLabel("by " + client.db.nameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02")))
		}
	}
	return client.threadPane
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}
