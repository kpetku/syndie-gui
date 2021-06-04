package ui

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/fetcher"
)

func (client *UI) fetchFromLocalFile() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, client.window)
			return
		}
		if reader == nil {
			return
		}
		err = fetchMessage("file", reader.URI().Path(), false)
		if err != nil {
			dg := dialog.NewError(err, client.window)
			dg.Show()
			return
		}
		dg := dialog.NewInformation("Open file", "Message imported", client.window)
		dg.Show()
		client.Rehash()
	}, client.window)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".syndie"}))
	fd.Show()
}

func (client *UI) fetchFromLocalFolder() {
	fd := dialog.NewFolderOpen(func(file fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, client.window)
			return
		}
		if file == nil {
			return
		}
		hbox := container.NewHBox()
		pb := widget.NewProgressBarInfinite()
		hbox.Add(widget.NewLabel("Fetching from " + file.Path()))
		hbox.Add(pb)
		pb.Show()
		progress := dialog.NewCustom("Fetching", "Cancel fetch", hbox, client.window)
		// TODO: allow fetch to be cancelled by passing a chan
		progress.Show()
		err = fetchMessage("directory", file.Path(), false)
		if err != nil {
			dg := dialog.NewError(err, client.window)
			dg.Show()
			return
		}
		progress.Hide()
		dg := dialog.NewInformation("Open directory", "Messages imported", client.window)
		dg.Show()
		client.Rehash()
	}, client.window)
	fd.Show()
}

func (client *UI) fetchFromArchiveServer(anonOnly bool) {
	content := container.NewVBox()
	selectedFetchArchive := widget.NewEntry()
	content.Add(widget.NewLabel("Press fetch to pull messages from the URL below"))
	selectedFetchArchive.SetPlaceHolder("http://localhost:8080/")
	cc := dialog.NewCustomConfirm("Fetch", "Fetch", "Cancel", content, func(fetch bool) {
		if selectedFetchArchive.Text == "" {
			selectedFetchArchive.Text = selectedFetchArchive.PlaceHolder
		}
		selectedFetchArchive.Text = "http://" + strings.TrimPrefix(selectedFetchArchive.Text, "http://")
		hbox := container.NewHBox()
		pb := widget.NewProgressBarInfinite()
		hbox.Add(widget.NewLabel("Fetching from " + selectedFetchArchive.Text))
		hbox.Add(pb)
		progress := dialog.NewCustom("Fetching", "Cancel fetch", hbox, client.window)
		progress.Show()
		err := fetchMessage("remoteURL", selectedFetchArchive.Text, anonOnly)
		if err != nil {
			progress.Hide()
			de := dialog.NewError(err, client.window)
			de.Show()
		}
		client.db.Reload()
		progress.Hide()
		client.Rehash()
	}, client.window)
	content.Add(selectedFetchArchive)
	cc.Show()
}

func fetchMessage(selectedFetchMethod string, location string, anonOnly bool) (err error) {
	switch selectedFetchMethod {
	case "file":
		f := fetcher.Fetcher{}
		err := f.LocalFile(location)
		if err != nil {
			return err
		}
	case "directory":
		f := fetcher.Fetcher{}
		err := f.LocalDir(location)
		if err != nil {
			return err
		}
	case "remoteURL":
		var f *fetcher.Fetcher
		dir, err := ioutil.TempDir("", "syndie")
		if err != nil {
			log.Fatalf("Unable to create a temporary directory: %s", err)
		}
		defer os.RemoveAll(dir)
		if anonOnly {
			f = fetcher.NewAnonOnly(location, dir, 60, 50)
		} else {
			f = fetcher.New(location, dir, 60, 50)
		}
		err = f.RemoteFetch()
		if err != nil {
			return err
		}
	}
	return nil
}
