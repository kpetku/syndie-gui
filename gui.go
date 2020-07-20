package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type GUI struct {
	db     *database
	window fyne.Window

	channelPane     *widget.Box
	threadPane      *widget.Box
	channelList     *widget.ScrollContainer
	messageList     *widget.ScrollContainer
	contentArea     *widget.Label
	selectedChannel string
	selectedMessage string
}

func NewGUI() *GUI {
	return new(GUI)
}

func (client *GUI) Start(path string) {
	client.db = NewDatabase()
	client.db.openDB(path)
	client.db.loadChannels()
	client.db.loadMessages()
	client.db.loadAvatars()

	a := app.New()

	client.window = a.NewWindow("Syndie GUI")
	client.contentArea = widget.NewLabel("This is where the message content goes")

	client.repaint()
	client.window.ShowAndRun()
}

func (client *GUI) repaint() {
	client.channelList = widget.NewVScrollContainer(client.renderChannelList())
	client.messageList = widget.NewVScrollContainer(client.renderThreadList())
	client.window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewHSplitContainer(client.channelList, client.messageList), client.contentArea))
}
