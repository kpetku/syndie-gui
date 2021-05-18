package main

import (
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

const version = "v0.0.2"

// GUI contains various GUI configuration options
type GUI struct {
	db     *database
	window fyne.Window

	channelPane *fyne.Container
	threadPane  *fyne.Container
	contentPane *fyne.Container
	channelList *container.Scroll
	messageList *container.Scroll
	contentArea *container.Scroll

	pagination      int
	selectedChannel string
	channelNeedle   int
	selectedMessage int
}

// NewGUI creates a new GUI
func NewGUI() *GUI {
	return new(GUI)
}

// Start launches a new syndie-gui application
func (client *GUI) Start(path string) {
	sanityCheckStartupDir(path)
	client.db = newDatabase()
	client.db.openDB(path + "/db/bolt.db")
	client.db.reload()

	a := app.New()

	client.window = a.NewWindow("Syndie" + version)
	client.renderMainMenu()
	client.applyOptions()

	client.repaint()
	client.window.Resize(fyne.NewSize(800, 600))
	client.window.ShowAndRun()
}

func (client *GUI) repaint() {
	client.channelList = container.NewVScroll(client.renderChannelList())
	client.messageList = container.NewVScroll(client.renderThreadList(client.channelNeedle))
	client.contentArea = container.NewVScroll(client.renderContentArea())
	client.window.SetContent(container.New(layout.NewGridLayout(3), client.channelList, client.messageList, client.contentArea))
}

func (client *GUI) applyOptions() {
	client.pagination = 25
}

func sanityCheckStartupDir(path string) {
	var err error
	var isWindows bool
	if runtime.GOOS == "windows" {
		isWindows = true
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		if isWindows {
			os.Mkdir(path, 0777)
		} else {
			os.Mkdir(path, 0700)
		}
	}
	_, err = os.Stat(path + "/db/")
	if os.IsNotExist(err) {
		if isWindows {
			os.Mkdir(path+"/db/", 0777)
		} else {
			os.Mkdir(path+"/db/", 0700)
		}
	}
}
