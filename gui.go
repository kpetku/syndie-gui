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
	contentPane     *widget.Box
	channelList     *widget.ScrollContainer
	messageList     *widget.ScrollContainer
	contentArea     *widget.ScrollContainer
	selectedChannel string
	selectedMessage int
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

	client.repaint()
	client.window.ShowAndRun()
}

func (client *GUI) repaint() {
	client.channelList = widget.NewVScrollContainer(client.renderChannelList())
	client.messageList = widget.NewVScrollContainer(client.renderThreadList())
	client.contentArea = widget.NewVScrollContainer(client.renderContentArea())
	client.window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(3), client.channelList, client.messageList, client.contentArea))
}
