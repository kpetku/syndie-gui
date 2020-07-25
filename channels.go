package main

import (
	"image/color"
	"strconv"
	"time"

	"fyne.io/fyne/canvas"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/kpetku/syndie-core/data"
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
		hbox := widget.NewHBox()

		icon := canvas.NewImageFromImage(client.db.getAvatar(channel.IdentHash))
		icon.SetMinSize(fyne.NewSize(32, 32))

		rw := new(tappableLabel)
		rw.SetText(channel.Name + " " + shortIdent(channel.IdentHash))
		rw.chanID = channel.IdentHash

		rw.selectedChannel = make(chan string)
		go func() {
			for click := range rw.selectedChannel {
				c := click
				client.selectedChannel = c
				client.channelNeedle = 0
				break
			}
			go client.repaint()
		}()

		hbox.Append(icon)
		hbox.Append(rw)
		client.channelPane.Append(hbox)
	}
	return client.channelPane
}

func (client *GUI) renderThreadList(needle int) fyne.CanvasObject {
	client.threadPane = widget.NewVBox()
	if client.selectedChannel == "" {
		client.threadPane.Append(widget.NewLabel("Select a channel from the menu to the left to get started"))
	} else {
		for num, msg := range client.db.chanList[client.selectedChannel] {
			if num <= client.channelNeedle || num < client.pagination {
				currentMessage := msg
				// TODO: Move this into it's own custom widget
				first := new(tappableLabel)
				first.msg = &currentMessage
				first.SetText(msg.Subject)

				first.selectedMessage = make(chan int)
				go func() {
					for click := range first.selectedMessage {
						c := click
						client.selectedMessage = c
						break
					}
					go client.repaint()
				}()
				client.threadPane.Append(first)

				second := new(tappableLabel)
				second.msg = &currentMessage
				date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
				second.SetText("by " + client.db.nameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02"))
				client.threadPane.Append(second)
			}
			if num == client.channelNeedle {
				if num <= client.pagination {
					continue
				}
				client.threadPane.Append(widget.NewButton("Show more messages", func() {
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
				client.threadPane.Append(widget.NewButton("Show more messages", func() {
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
	client.contentPane = widget.NewVBox()
	if client.selectedChannel == "" {
		return client.contentPane
	} else {
		if client.selectedMessage != 0 {
			var currentMessage data.Message
			for _, msg := range client.db.chanList[client.selectedChannel] {
				if msg.ID == client.selectedMessage {
					currentMessage = msg
					break
				}
			}
			if currentMessage.Subject != "" {
				subject := widget.NewLabel("Subject: " + currentMessage.Subject)
				subject.Wrapping = fyne.TextWrapBreak
				client.contentPane.Append(subject)
				client.contentPane.Append(canvas.NewLine(color.White))
				if len(currentMessage.Raw.Page) > 0 {
					for num, p := range currentMessage.Raw.Page[:1] {
						if num >= 0 {
							client.contentPane.Append(widget.NewLabel("Page: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(currentMessage.Raw.Page)-1)))
							page := widget.NewLabel(p.Data)
							page.Wrapping = fyne.TextWrapBreak
							client.contentPane.Append(page)
							client.contentPane.Append(canvas.NewLine(color.White))
						}
					}
				}
				if len(currentMessage.Raw.Attachment) > 0 {
					for num, a := range currentMessage.Raw.Attachment {
						if num >= 0 {
							a1 := widget.NewLabel("Attachment: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(currentMessage.Raw.Attachment)) + " Name: " + a.Name)
							a1.Wrapping = fyne.TextWrapBreak
							a2 := widget.NewLabel("Type: " + a.ContentType + " Description: " + a.Description)
							a2.Wrapping = fyne.TextWrapBreak
							client.contentPane.Append(a1)
							client.contentPane.Append(a2)
							adata := a.Data
							image, err := renderImage(imageExtFromName(a.ContentType), adata)
							if err != nil {
								client.contentPane.Append(widget.NewLabel("Unable to display preview"))
							} else {
								i := canvas.NewImageFromImage(image)
								i.FillMode = canvas.ImageFillContain
								i.SetMinSize(fyne.NewSize(fyne.Min(image.Bounds().Dx(), client.contentArea.Size().Width), image.Bounds().Dy()))
								client.contentPane.Append(i)
							}
							client.contentPane.Append(canvas.NewLine(color.White))
						}
					}
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
