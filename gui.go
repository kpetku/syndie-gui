package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/kpetku/syndie-core/fetcher"
)

// GUI contains various GUI configuration options
type GUI struct {
	db     *database
	window fyne.Window

	channelPane *widget.Box
	threadPane  *widget.Box
	contentPane *widget.Box
	channelList *widget.ScrollContainer
	messageList *widget.ScrollContainer
	contentArea *widget.ScrollContainer

	pagination           int
	selectedChannel      string
	channelNeedle        int
	selectedMessage      int
	selectedFetchArchive *widget.Entry
	selectedFetchMethod  string
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

	client.window = a.NewWindow("Syndie GUI")
	client.loadMainMenu()
	client.applyOptions()

	client.repaint()
	client.window.Resize(fyne.NewSize(800, 600))
	client.window.ShowAndRun()
}

func (client *GUI) repaint() {
	client.channelList = widget.NewVScrollContainer(client.renderChannelList())
	client.messageList = widget.NewVScrollContainer(client.renderThreadList(client.channelNeedle))
	client.contentArea = widget.NewVScrollContainer(client.renderContentArea())
	client.window.SetContent(fyne.NewContainerWithLayout(layout.NewGridLayout(3), client.channelList, client.messageList, client.contentArea))
}

func (client *GUI) loadMainMenu() {
	main := fyne.NewMenu("File",
		fyne.NewMenuItem("Fetch from URL", func() {
			client.selectedFetchMethod = "URL"
			client.fetchFromArchiveServer()
		}),
		fyne.NewMenuItem("Fetch from directory", func() {
			client.selectedFetchMethod = "directory"
			client.fetchFromArchiveServer()
		}),
	)
	client.window.SetMainMenu(fyne.NewMainMenu(main))
}

func (client *GUI) fetchFromArchiveServer() {
	content := widget.NewVBox()
	client.selectedFetchArchive = widget.NewEntry()
	content.Append(widget.NewLabel("Press fetch to pull messages from the " + client.selectedFetchMethod + " below"))
	if client.selectedFetchMethod == "URL" {
		client.selectedFetchArchive.SetPlaceHolder("http://localhost:8080/")
	}
	if client.selectedFetchMethod == "directory" {
		client.selectedFetchArchive.SetPlaceHolder("~/.syndie/archive")
	}
	dialog.NewCustomConfirm("Fetch", "Fetch", "Cancel", content, client.fetch, client.window)
	content.Append(client.selectedFetchArchive)
}

func (client *GUI) fetch(fetch bool) {
	if client.selectedFetchArchive.Text == "" {
		client.selectedFetchArchive.Text = client.selectedFetchArchive.PlaceHolder
	}
	if client.selectedFetchMethod == "URL" {
		client.selectedFetchArchive.Text = "http://" + strings.TrimLeft(client.selectedFetchArchive.Text, "http://")
	}
	dir, err := ioutil.TempDir("", "syndie")
	if err != nil {
		log.Fatalf("Unable to create a temporary directory: %s", err)
	}
	defer os.RemoveAll(dir)
	if fetch {
		progress := dialog.NewProgressInfinite("Fetching", "Fetching from "+client.selectedFetchArchive.Text, client.window)
		f := fetcher.New(client.selectedFetchArchive.Text, dir, 60, 50)
		progress.Show()
		if client.selectedFetchMethod == "URL" {
			err := f.RemoteFetch()
			if err != nil {
				progress.Hide()
				dialog.NewError(err, client.window)
			}
		}
		if client.selectedFetchMethod == "directory" {
			client.selectedFetchArchive.Text = strings.TrimRight(client.selectedFetchArchive.Text, "/") + "/"
			err := f.LocalDir(client.selectedFetchArchive.Text)
			log.Printf("Checking: %s", client.selectedFetchArchive.Text)
			if err != nil {
				progress.Hide()
				dialog.NewError(err, client.window)
			}
		}
		client.db.reload()
		progress.Hide()
		client.repaint()
	}
}

func (client *GUI) applyOptions() {
	client.pagination = 25
}

func sanityCheckStartupDir(path string) {
	var err error
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		os.Mkdir(path, 0777)
		_, err = os.Stat(path + "/db/")
		if os.IsNotExist(err) {
			os.Mkdir(path+"/db/", 0777)
		}
	}
}
