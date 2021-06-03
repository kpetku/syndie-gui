package ui

import (
	"fyne.io/fyne/v2"
)

func (client *UI) renderMainMenu() *fyne.MainMenu {
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
	return fyne.NewMainMenu(main)
}
