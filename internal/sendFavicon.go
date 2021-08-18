package internal

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strconv"
)

func sendFavicon(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	m := image.NewRGBA(image.Rect(0, 0, 16, 16))
	clr := color.RGBA{0, 255, 255, 255}
	draw.Draw(m, m.Bounds(), image.NewUniform(clr), image.Point{}, draw.Src)
	jpeg.Encode(buf, m, nil)
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf.Bytes())))
	w.Write(buf.Bytes())
}
