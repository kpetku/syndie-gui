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
	"github.com/kpetku/syndie-gui/database"
)

// GUI contains various GUI configuration options
type GUI struct {
	db     *database.Database
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
	client.db = database.New()
	client.db.Open(path + "/db/bolt.db")
	client.db.Reload()

	a := app.New()

	client.window = a.NewWindow("Syndie " + version)
	client.renderMainMenu()
	client.applyOptions()
	client.preloadAvatarCache()
	client.window.Resize(fyne.NewSize(800, 600))
	client.window.SetContent(client.renderFeedView())
	client.window.ShowAndRun()
}

// Rehash reloads the database, reloads the avatar cache, and repaints the main window
func (client *GUI) Rehash() {
	client.db.Reload()
	client.preloadAvatarCache()
	client.repaintMainWindow()
}

func (client *GUI) preloadAvatarCache() {
	client.avatarCache = make(map[string]*canvas.Image)
	for _, channel := range client.db.Channels {
		client.avatarCache[channel.IdentHash] = canvas.NewImageFromImage(client.db.GetAvatar(channel.IdentHash))
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

func (client GUI) renderNavBar(highlighted string) *fyne.Container {
	feedButton := widget.NewButton("Feed", func() { client.window.SetContent(client.renderFeedView()) })
	latestButton := widget.NewButton("Latest", func() { client.window.SetContent(client.renderLatestView()) })
	forYouButton := widget.NewButton("For you", func() {})
	followingButton := widget.NewButton("Following", func() {})

	navBar := container.NewGridWithColumns(4)
	navBar.Add(feedButton)
	navBar.Add(forYouButton)
	navBar.Add(latestButton)
	navBar.Add(followingButton)
	switch highlighted {
	case "feed":
		feedButton.Importance = widget.HighImportance
	case "latest":
		latestButton.Importance = widget.HighImportance
	case "foryou":
		forYouButton.Importance = widget.HighImportance
	case "following":
		followingButton.Importance = widget.HighImportance
	}
	return navBar
}
