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
		fyne.NewMenuItem("Fetch from using I2P or Tor (anonymous)", func() {
			client.fetchFromArchiveServer(true)
		}),
		fyne.NewMenuItem("Fetch from the Internet (non-anonymous)", func() {
			client.fetchFromArchiveServer(false)
		}),
	)
	return fyne.NewMainMenu(main)
}
