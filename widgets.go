package main

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type ChannelWidget struct {
	widget.BaseWidget
	Name            string
	Icon            *canvas.Image
	Text            *widget.Label
	SelectedChannel chan string
	Highlight       bool
}

type ChannelWidgetRenderer struct {
	icon      *canvas.Image
	text      *widget.Label
	name      *widget.Label
	highlight bool
	objects   []fyne.CanvasObject
}

func (rw *ChannelWidget) Tapped(_ *fyne.PointEvent) {
	rw.SelectedChannel <- rw.Text.Text
}

func (rw *ChannelWidget) TappedSecondary(_ *fyne.PointEvent) {
}

func (rw *ChannelWidget) CreateRenderer() fyne.WidgetRenderer {
	rw.ExtendBaseWidget(rw)
	n := widget.NewLabelWithStyle(rw.Name, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	rwr := &ChannelWidgetRenderer{}
	rwr.icon = rw.Icon
	rwr.name = n
	rwr.text = rw.Text
	rwr.highlight = rw.Highlight
	rwr.objects = []fyne.CanvasObject{rw.Icon, n, rw.Text}
	return rwr
}

func (rwr *ChannelWidgetRenderer) Layout(size fyne.Size) {
	rwr.name.Move(fyne.NewPos(0, 0))
	rwr.text.Move(fyne.NewPos(0, 25))
	// Move everything over 50 for the icon
	rwr.name.Move(fyne.NewPos(50+rwr.name.Position().X, rwr.name.Position().Y))
	rwr.text.Move(fyne.NewPos(50+rwr.text.Position().X, rwr.text.Position().Y))
	// Display the icon
	rwr.icon.Move(fyne.NewPos(0, 0))
	rwr.icon.Resize(fyne.NewSize(50, 50))
}

func (rwr *ChannelWidgetRenderer) Objects() []fyne.CanvasObject {
	return rwr.objects
}

func (rw *ChannelWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 50)
}

func (rwr *ChannelWidgetRenderer) BackgroundColor() color.Color {
	if rwr.highlight {
		return theme.HoverColor()
	}
	return theme.BackgroundColor()
}

func (rw *ChannelWidgetRenderer) Destroy() {
}

func (rwr *ChannelWidgetRenderer) Refresh() {
	rwr.objects = append(rwr.objects, rwr.icon)
	rwr.objects = append(rwr.objects, rwr.name)
	rwr.objects = append(rwr.objects, rwr.text)
}
