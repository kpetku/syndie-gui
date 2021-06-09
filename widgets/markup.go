package widgets

import (
	"bytes"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"golang.org/x/net/html"
)

func NewMarkup(s string) (*fyne.Container, error) {
	var sb strings.Builder
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
				sb.Reset()
				vbox.Add(NewLabel("Hyperlink placeholder: " + renderNode(n)))
				return
			case "img":
				sb.Reset()
				vbox.Add(NewLabel("Image placeholder: " + renderNode(n)))
				return
			case "br":
				vbox.Add(layout.NewSpacer())
				sb.Reset()
				return
			default:
			}
		case html.TextNode:
			if renderNode(n) != "" || sb.Len() > 0 {
				sb.WriteString(renderNode(n))
				if html.UnescapeString(sb.String()) != "" {
					vbox.Add(NewLabel(html.UnescapeString(sb.String())))
				}
			}
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
