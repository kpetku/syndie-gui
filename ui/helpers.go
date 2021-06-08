package ui

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/kpetku/syndie-core/data"
	"github.com/kpetku/syndie-gui/widgets"
)

func newCenteredContainer(l *fyne.Container) *fyne.Container {
	if !fyne.CurrentDevice().IsMobile() {
		// Center the layout on desktop
		return container.New(layout.NewGridLayoutWithColumns(1), l)
	}
	return l
}

func (client UI) msgToCard(msg data.Message) *widget.Card {
	date := time.Unix(0, int64(msg.ID)*int64(time.Millisecond))
	text := "by " + client.db.NameFromChanIdentHash(msg.Author) + " " + shortIdent(msg.Author) + " on " + date.Format("2006-01-02")
	vbox := container.NewVBox()

	hbox := container.NewHBox()
	icon := client.avatarCache[msg.Author]
	if icon != nil {
		hbox.Add(icon)
	}
	hbox.Add(widgets.NewLabel(text))

	vbox.Add(hbox)
	if len(msg.Raw.Page) > 0 {
		for num, p := range msg.Raw.Page[:1] {
			if num >= 0 {
				vbox.Add(widgets.NewLabel("Page: " + strconv.Itoa(num+1) + "/" + strconv.Itoa(len(msg.Raw.Page)-1)))
				vbox.Add(widget.NewSeparator())
				vbox.Add(widgets.NewLabel(p.Data))
			}
		}
	}
	return widget.NewCard(msg.Subject, "", vbox)
}

func shortIdent(i string) string {
	if len(i) > 6 {
		return "[" + i[0:6] + "]"
	}
	return ""
}

func renderImage(ext string, data []byte) (image.Image, error) {
	var image image.Image
	var err error
	switch ext {
	case "png":
		image, err = png.Decode(bytes.NewReader(data))
	case "jpeg":
		image, err = jpeg.Decode(bytes.NewReader(data))
	case "gif":
		image, err = gif.Decode(bytes.NewReader(data))
	default:
		image, err = jpeg.Decode(bytes.NewReader(data))
	}
	return image, err
}

func imageExtFromName(s string) string {
	if strings.Contains(s, "/") {
		switch strings.Split(s, "/")[1] {
		case "gif":
			return "gif"
		case "png":
			return "png"
		case "jpg", "jpeg":
			return "jpeg"
		}
	}
	return ""
}