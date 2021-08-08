package imgtxtcolor

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
)

func (p *stParam) ToGif() {
	var images []*image.Paletted
	var delays []int
	var disposals []byte
	//var palette2 color.Palette = palette.Plan9 // TODO 256 цветов https://pkg.go.dev/image/color/palette@go1.16.6
	// в Plan9 нет прозрачного цвета
	var palette2 color.Palette = palette.WebSafe   // 216 цветов
	palette2 = append(palette2, image.Transparent) // добавляем еще прозрачный
	for i := 0; i < len(p.allCanvas); i++ {

		img := p.allCanvas[i].Img
		bounds := img.Bounds()

		dst := image.NewPaletted(bounds, palette2)
		draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

		images = append(images, dst)
		delays = append(delays, p.allCanvas[i].GifDelay)
		disposals = append(disposals, gif.DisposalBackground)
	}

	f, err := os.OpenFile(p.opt.GifFileName, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()
	err = gif.EncodeAll(f, &gif.GIF{
		Image:    images,
		Delay:    delays,
		Disposal: disposals, //без этго не стирает предыдущие картинки gif.DisposalBackground
	})
	if err != nil {
		log.Println(err.Error())
	}
	//log.Printf("%+v", palette2)
}
