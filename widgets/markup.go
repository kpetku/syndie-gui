package widgets

import (
	"bytes"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"golang.org/x/net/html"
)

func NewMarkup(s string) (*fyne.Container, error) {
	vbox := container.NewVBox()
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return &fyne.Container{}, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		switch n.Type {
		case html.ElementNode:
			switch n.Data {
			case "quote":
			case "a":
			case "img":
			//	log.Printf("Image placeholder: %s", renderNode(n))
			case "br":
			default:
			}
		case html.TextNode:
			// log.Printf("%s", renderNode(n))
		case html.DocumentNode:
		default:
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return vbox, nil
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
