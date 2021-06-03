package ui

import (
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"

	"github.com/kpetku/syndie-gui/database"
)

const version = "0.0.2"

type UI struct {
	app    fyne.App
	db     *database.Database
	window fyne.Window

	avatarCache map[string]*canvas.Image

	pagination      int
	selectedChannel string
	channelNeedle   int
	threadPane      *fyne.Container
}

// NewUI creates a new UI
func NewUI() *UI {
	u := new(UI)
	u.avatarCache = make(map[string]*canvas.Image)
	return u
}

// Start launches a new syndie-UI application
func (client *UI) Start(path string) {
	sanityCheckStartupDir(path)
	client.db = database.New()
	client.db.Open(path + "/db/bolt.db")
	client.db.Reload()

	client.app = app.New()
	client.window = client.app.NewWindow("Syndie " + version)
	client.applyOptions()
	client.preloadAvatarCache()
	client.window.Resize(fyne.NewSize(800, 600))

	client.window.SetMainMenu(client.renderMainMenu())
	client.SetContent(client.renderFeedView())
	client.window.ShowAndRun()
}

// Rehash reloads the database, reloads the avatar cache, and repaints the main window
func (client *UI) Rehash() {
	client.db.Reload()
	client.preloadAvatarCache()
	client.repaintMainWindow()
}

func (client *UI) preloadAvatarCache() {
	for _, channel := range client.db.Channels {
		client.avatarCache[channel.IdentHash] = canvas.NewImageFromImage(client.db.GetAvatar(channel.IdentHash))
		client.avatarCache[channel.IdentHash].SetMinSize(fyne.NewSize(32, 32))
		client.avatarCache[channel.IdentHash].FillMode = canvas.ImageFillContain
	}
}

func (client *UI) repaintMainWindow() {
	if client.selectedChannel != "" {
		client.window.SetContent(container.NewScroll(newCenteredContainer(client.renderThreadListWithMenu(0))))
		return
	}
	client.window.SetContent(client.renderFeedView())
}

func (client *UI) SetContent(o fyne.CanvasObject) {
	client.window.SetContent(o)
}

func (client *UI) applyOptions() {
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
