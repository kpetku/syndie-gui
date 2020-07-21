package main

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"
)

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
	switch filepath.Ext(strings.Replace(s, "/", ".", 0)) {
	case "gif":
		return "gif"
	case "png":
		return "png"
	case "jpg", "jpeg":
		return "jpeg"
	}
	return ""
}
