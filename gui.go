package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/fetcher"
)

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
	client.channelList = container.NewVScroll(client.renderChannelList())
	client.messageList = container.NewVScroll(client.renderThreadList(client.channelNeedle))
	client.contentArea = container.NewVScroll(client.renderContentArea())
	client.window.SetContent(container.New(layout.NewGridLayout(3), client.channelList, client.messageList, client.contentArea))
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
	content := container.NewVBox()
	client.selectedFetchArchive = widget.NewEntry()
	content.Add(widget.NewLabel("Press fetch to pull messages from the " + client.selectedFetchMethod + " below"))
	if client.selectedFetchMethod == "URL" {
		client.selectedFetchArchive.SetPlaceHolder("http://localhost:8080/")
	}
	if client.selectedFetchMethod == "directory" {
		client.selectedFetchArchive.SetPlaceHolder("~/.syndie/archive")
	}
	cc := dialog.NewCustomConfirm("Fetch", "Fetch", "Cancel", content, client.fetch, client.window)
	content.Add(client.selectedFetchArchive)
	cc.Show()
}

func (client *GUI) fetch(fetch bool) {
	if client.selectedFetchArchive.Text == "" {
		client.selectedFetchArchive.Text = client.selectedFetchArchive.PlaceHolder
	}
	if client.selectedFetchMethod == "URL" {
		client.selectedFetchArchive.Text = "http://" + strings.TrimPrefix(client.selectedFetchArchive.Text, "http://")
	}
	dir, err := ioutil.TempDir("", "syndie")
	if err != nil {
		log.Fatalf("Unable to create a temporary directory: %s", err)
	}
	defer os.RemoveAll(dir)
	if fetch {
		hbox := container.NewHBox()
		pb := widget.NewProgressBarInfinite()
		hbox.Add(widget.NewLabel("Fetching from " + client.selectedFetchArchive.Text))
		hbox.Add(pb)
		progress := dialog.NewCustom("Fetching", "Cancel fetch", hbox, client.window)
		f := fetcher.New(client.selectedFetchArchive.Text, dir, 60, 50)
		progress.Show()
		if client.selectedFetchMethod == "URL" {
			err := f.RemoteFetch()
			if err != nil {
				progress.Hide()
				de := dialog.NewError(err, client.window)
				de.Show()
			}
		}
		if client.selectedFetchMethod == "directory" {
			client.selectedFetchArchive.Text = strings.TrimRight(client.selectedFetchArchive.Text, "/") + "/"
			err := f.LocalDir(client.selectedFetchArchive.Text)
			log.Printf("Checking: %s", client.selectedFetchArchive.Text)
			if err != nil {
				progress.Hide()
				de := dialog.NewError(err, client.window)
				de.Show()
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
