package widgets

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/kpetku/libsyndie/syndieutil"
	"github.com/kpetku/syndie-core/data"
	"golang.org/x/net/html"
)

func NewMarkup(msg data.Message, s string) (*fyne.Container, error) {
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
				i, err := parseImg(msg, renderNode(n))
				if err != nil {
					log.Printf("parseImg error: %s", err)
					return
				}
				vbox.Add(i)
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
				if strings.TrimSpace(html.UnescapeString(sb.String())) != "" {
					vbox.Add(NewLabel(html.UnescapeString(sb.String())))
				}
				sb.Reset()
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

func parseImg(msg data.Message, s string) (*canvas.Image, error) {
	var i *canvas.Image
	tok := html.NewTokenizer(strings.NewReader(s))
	for {
		tokType := tok.Next()
		if tokType == html.ErrorToken {
			err := tok.Err()
			if err == io.EOF {
				break
			}
			return i, errors.New("error tokenizing html")
		}
		if tokType == html.SelfClosingTagToken {
			u := tok.Token().Attr[0].Val
			out := syndieutil.URI{}
			out.Marshall(u)
			if out.Channel == "" && out.Attachment > 0 {
				size, ext, err := image.DecodeConfig(bytes.NewReader(msg.Raw.Attachment[out.Attachment-1].Data))
				if err != nil {
					log.Printf("Error rendering image: %s", err)
					return &canvas.Image{}, err
				}
				raw, err := renderImage(ext, msg.Raw.Attachment[out.Attachment-1].Data)
				if err != nil {
					log.Printf("Error rendering image: %s", err)
					return &canvas.Image{}, err
				}
				i = canvas.NewImageFromImage(raw)
				i.FillMode = canvas.ImageFillContain
				i.SetMinSize(fyne.NewSize(float32(size.Width), float32(size.Height)))
				return i, nil
			}
			log.Printf("External image from attachment unimplemented: %s", u)
			return &canvas.Image{}, nil
		}
	}
	return i, nil
}

func renderImage(ext string, data []byte) (image.Image, error) {
	var image image.Image
	var err error
	switch ext {
	case "png":
		image, err = png.Decode(bytes.NewReader(data))
	case "jpeg", "jpg":
		image, err = jpeg.Decode(bytes.NewReader(data))
	case "gif":
		image, err = gif.Decode(bytes.NewReader(data))
	default:
		image, err = jpeg.Decode(bytes.NewReader(data))
	}
	return image, err
}
