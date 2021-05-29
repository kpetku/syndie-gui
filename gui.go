package main

import (
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const version = "v0.0.2"

// GUI contains various GUI configuration options
type GUI struct {
	db     *database
	window fyne.Window

	threadPane  *fyne.Container
	avatarCache map[string]*canvas.Image

	pagination      int
	selectedChannel string
	channelNeedle   int
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

	client.window = a.NewWindow("Syndie " + version)
	client.renderMainMenu()
	client.applyOptions()
	client.preloadAvatarCache()
	client.window.Resize(fyne.NewSize(800, 600))
	client.window.SetContent(client.renderFeedView())
	client.window.ShowAndRun()
}

func (client *GUI) preloadAvatarCache() {
	client.avatarCache = make(map[string]*canvas.Image)
	for _, channel := range client.db.Channels {
		client.avatarCache[channel.IdentHash] = canvas.NewImageFromImage(client.db.getAvatar(channel.IdentHash))
		client.avatarCache[channel.IdentHash].SetMinSize(fyne.NewSize(32, 32))
		client.avatarCache[channel.IdentHash].FillMode = canvas.ImageFillContain
	}
}

func (client *GUI) repaintMainWindow() {
	if client.selectedChannel != "" {
		client.window.SetContent(container.NewScroll(newCenteredContainer(client.renderThreadListWithMenu(0))))
		return
	}
	client.window.SetContent(client.renderFeedView())
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

func newCenteredContainer(l *fyne.Container) *fyne.Container {
	empty := widget.NewLabel("")
	if !fyne.CurrentDevice().IsMobile() {
		// Center the layout on desktop
		return container.New(layout.NewGridLayout(3), empty, l, empty)
	}
	return l
}
