package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

type GUI struct {
	db *database
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

	w := a.NewWindow("Syndie GUI")
	rightSideBar := widget.NewVScrollContainer(client.renderChannelList())
	w.SetContent(rightSideBar)
	w.ShowAndRun()
}
