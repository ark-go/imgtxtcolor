package imgtxtcolor

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

func ToGif(param *stParam, fileName string) {
	var images []*image.Paletted
	var delays []int
	var disposals []byte
	//var palette2 color.Palette = palette.Plan9 // TODO 256 цветов https://pkg.go.dev/image/color/palette@go1.16.6
	// в Plan9 нет прозрачного цвета
	var palette2 color.Palette = palette.WebSafe   // 216 цветов
	palette2 = append(palette2, image.Transparent) // добавляем еще прозрачный
	for i := 0; i < len(param.allImages); i++ {

		img := param.allImages[i]
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
	err = gif.EncodeAll(f, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals,
	})
	if err != nil {
		log.Println(err.Error())
	}
	//log.Printf("%+v", palette2)
}
