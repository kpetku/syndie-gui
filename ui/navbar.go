package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (client *UI) NewNavbar(highlighted string) *fyne.Container {
	// TODO: Allow the GUI to call SetContent on a window
	feedButton := widget.NewButton("Feed", func() {
		client.window.SetContent(client.renderFeedView())
	})
	latestButton := widget.NewButton("Latest", func() {
		client.window.SetContent(client.renderLatestView())
	})
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
