package imgtxtcolor

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"

	"github.com/fogleman/gg"
)

func ToGif(param *stParam, fileName string) {
	const width, height = 300, 300
	var images []*image.Paletted
	var delays []int
	var disposals []byte
	dc := gg.NewContext(width, height)

	for i := 0; i < len(param.allImages); i++ {
		dc.SetRGBA(1, 1, 1, 0)
		dc.Clear()

		img := param.allImages[i]

		var palette2 color.Palette = palette.Plan9 // TODO 256 цветов https://pkg.go.dev/image/color/palette@go1.16.6
		bounds := img.Bounds()

		dst := image.NewPaletted(bounds, palette2)

		draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

		images = append(images, dst)
		delays = append(delays, param.opt.GifDelay)
		disposals = append(disposals, gif.DisposalBackground)
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	gif.EncodeAll(f, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals,
	})
}
