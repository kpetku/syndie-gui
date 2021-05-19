package main

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/kpetku/syndie-core/data"
)

func (client *GUI) renderChannelList() fyne.CanvasObject {
	var hideEmptyChannels = true

	client.channelPane = container.NewVBox()
	for _, channel := range client.db.Channels {
		if hideEmptyChannels {
			if len(client.db.chanList[channel.IdentHash]) == 0 {
				continue
			}
		}
		hbox := container.NewHBox()

		icon := canvas.NewImageFromImage(client.db.getAvatar(channel.IdentHash))
		icon.SetMinSize(fyne.NewSize(32, 32))
		icon.FillMode = canvas.ImageFillContain

		rw := new(tappableLabel)
		rw.SetText(channel.Name + " " + shortIdent(channel.IdentHash) + " (?/" + strconv.Itoa(len(client.db.chanList[channel.IdentHash])) + ")")
		rw.chanID = channel.IdentHash

		rw.selectedChannel = make(chan string)
		go func() {
			for click := range rw.selectedChannel {
				c := click
				client.selectedChannel = c
				client.channelNeedle = 0
				client.selectedMessage = 0
				client.messageList.ScrollToTop()
				break
			}
			go client.repaint()
		}()

		hbox.Add(icon)
		hbox.Add(rw)
		client.channelPane.Add(hbox)
	}
	return client.channelPane
}

func (client *GUI) renderThreadList(needle int) fyne.CanvasObject {
	client.threadPane = container.NewVBox()
	if client.selectedChannel == "" {
		client.threadPane.Add(widget.NewLabel("Select a channel from the menu to the left to get started"))
	} else {
		for num, msg := range client.db.chanList[client.selectedChannel] {
			if num <= client.channelNeedle || num < client.pagination {
				currentMessage := msg
				first := newTappableLabel(msg.Subject)
				first.msg = &currentMessage
				first.selectedMessage = make(chan int)
				go func() {
					for click := range first.selectedMessage {
						c := click
						client.selectedMessage = c
						client.contentArea.ScrollToTop()
						break
					}
					go client.repaint()
				}()
				first.Alignment = fyne.TextAlignLeading
				client.threadPane.Add(first)
				date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
				second := widget.NewLabel("by " + client.db.nameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02"))
				second.TextStyle = fyne.TextStyle{Monospace: true}
				client.threadPane.Add(newLabelWithImage(canvas.NewImageFromImage(client.db.getAvatar(msg.Author)), second))
				client.threadPane.Add(canvas.NewLine(theme.ShadowColor()))
			}
			if num == client.channelNeedle {
				if num <= client.pagination {
					continue
				}
				client.threadPane.Add(widget.NewButton("Show more messages", func() {
					client.channelNeedle = needle + client.pagination
					client.renderThreadList(needle + client.pagination)
					client.repaint()
				}))
			}
		}
		if needle == client.channelNeedle {
			if len(client.db.chanList[client.selectedChannel])-1 <= client.pagination {
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
					client.repaint()
				}))
			}
		}
	}
	return client.threadPane
}

func (client *GUI) renderContentArea() fyne.CanvasObject {
	client.contentPane = container.NewVBox()
	if client.selectedChannel == "" {
		return client.contentPane
	}
	if client.selectedMessage != 0 {
		var currentMessage data.Message
		for _, msg := range client.db.chanList[client.selectedChannel] {
			if msg.ID == client.selectedMessage {
				currentMessage = msg
				break
			}
		}
		if currentMessage.Subject != "" {
			s := newLabel("Subject: " + currentMessage.Subject)
			s.TextStyle = fyne.TextStyle{Monospace: true}
			client.contentPane.Add(s)
		} else {
			s := newLabel("No subject")
			s.TextStyle = fyne.TextStyle{Monospace: true}
			client.contentPane.Add(s)
		}
		client.contentPane.Add(canvas.NewLine(theme.ShadowColor()))
		if len(currentMessage.Raw.Page) > 0 {
			for num, p := range currentMessage.Raw.Page[:1] {
				if num >= 0 {
					client.contentPane.Add(newLabel("Page: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(currentMessage.Raw.Page)-1)))
					client.contentPane.Add(newLabel(p.Data))
					client.contentPane.Add(canvas.NewLine(theme.ShadowColor()))
				}
			}
		}
		if len(currentMessage.Raw.Attachment) > 0 {
			for num, a := range currentMessage.Raw.Attachment {
				if num >= 0 {
					client.contentPane.Add(newLabel("Attachment: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(currentMessage.Raw.Attachment)) + " Name: " + a.Name))
					client.contentPane.Add(newLabel("Type: " + a.ContentType + " Description: " + a.Description))
					adata := a.Data
					image, err := renderImage(imageExtFromName(a.ContentType), adata)
					if err != nil {
						client.contentPane.Add(widget.NewLabel("Unable to display preview"))
					} else {
						i := canvas.NewImageFromImage(image)
						i.FillMode = canvas.ImageFillContain
						i.SetMinSize(fyne.NewSize(fyne.Min(float32(image.Bounds().Dx()), client.contentArea.Size().Width), float32(image.Bounds().Dy())))
						client.contentPane.Add(i)
					}
					client.contentPane.Add(canvas.NewLine(theme.ShadowColor()))
				}
			}
		}
	}
	return client.contentPane
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}
