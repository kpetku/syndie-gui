package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (client *UI) renderSettingsView() {
	content := container.NewVBox()
	content.Add(widget.NewCheckWithData("Hide empty feeds", binding.BindPreferenceBool("hideEmptyFeeds", fyne.CurrentApp().Preferences())))
	entries := container.New(layout.NewFormLayout())
	entries.Add(widget.NewLabel("Number of messages per page"))
	entries.Add(widget.NewEntryWithData(binding.BindPreferenceString("pagination", fyne.CurrentApp().Preferences())))
	content.Add(entries)
	dialog.ShowCustom("Settings", "Back", content, client.window)
}
