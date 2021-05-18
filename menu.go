package main

import (
	"fyne.io/fyne/v2"
)

func (client *GUI) renderMainMenu() {
	main := fyne.NewMenu("File",
		fyne.NewMenuItem("Open file", func() {
			client.fetchFromLocalFile()
		}),
		fyne.NewMenuItem("Open folder", func() {
			client.fetchFromLocalFolder()
		}),
		fyne.NewMenuItem("Fetch from URL", func() {
			client.fetchFromArchiveServer()
		}),
	)
	client.window.SetMainMenu(fyne.NewMainMenu(main))
}
