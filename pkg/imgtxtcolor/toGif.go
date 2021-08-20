package imgtxtcolor

import (
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"sync"
	"time"

	"github.com/ark-go/imgtxtcolor/pkg/quantize"
)

func (p *stParam) ToGif() {
	startTime := time.Now()
	count := len(p.allCanvas)
	images := make([]*image.Paletted, count)
	delays := make([]int, count)
	disposals := make([]byte, count)
	//var palette2 color.Palette = palette.Plan9 // TODO 256 цветов https://pkg.go.dev/image/color/palette@go1.16.6
	// в Plan9 нет прозрачного цвета
	// var palette2 color.Palette = palette.WebSafe   // 216 цветов
	// palette2 = append(palette2, image.Transparent) // добавляем еще прозрачный

	//var idx int
	var gifHeight int
	var gifWidth int
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < len(p.allCanvas); i++ {
		wg.Add(1)
		go createLayerGif(p, i, images, delays, disposals, &gifWidth, &gifHeight, &wg, &mu)

		// img := p.allCanvas[i].Img
		// bounds := img.Bounds()
		// //--------------------
		// q := quantize.MedianCutQuantizer{}
		// palette2 := q.Quantize(make([]color.Color, 0, 256), img)
		// //!------ TODO:  не понятно, мне в палитре нужен прозрачный цвет,
		// //!------       ищем наиболее близкий, и заменяем его прозрачным
		// idx = palette2.Index(color.RGBA{0, 0, 0, 0})
		// palette2[idx] = color.RGBA{0, 0, 0, 0}
		// //--------------------
		// dst := image.NewPaletted(bounds, palette2)
		// draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

		// //images = append(images, dst)
		// images[i] = dst
		// //delays = append(delays, p.allCanvas[i].GifDelay)
		// delays[i] = p.allCanvas[i].GifDelay
		// //	disposals = append(disposals, gif.DisposalBackground)
		// disposals[i] = gif.DisposalBackground
		// if img.Rect.Dx() > gifWidth {
		// 	gifWidth = img.Rect.Dx()
		// }
		// if img.Rect.Dy() > gifHeight {
		// 	gifHeight = img.Rect.Dy()
		// }

	}
	wg.Wait()
	log.Printf("Time [%v]: %v\n", "прошли For gif", time.Since(startTime))
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
		//	BackgroundIndex:// byte(idx),
		Config: image.Config{
			Width:  gifWidth,
			Height: gifHeight,
		},
	})
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Time [%v]: %v\n", "Записали gif", time.Since(startTime))
}

func createLayerGif(p *stParam, i int, images []*image.Paletted, delays []int, disposals []byte, gifWidth *int, gifHeight *int, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	img := p.allCanvas[i].Img
	bounds := img.Bounds()
	//--------------------
	q := quantize.MedianCutQuantizer{}
	palette2 := q.Quantize(make([]color.Color, 0, 256), img)
	//!------ TODO:  не понятно, мне в палитре нужен прозрачный цвет,
	//!------       ищем наиболее близкий, и заменяем его прозрачным
	idx := palette2.Index(color.RGBA{0, 0, 0, 0})
	palette2[idx] = color.RGBA{0, 0, 0, 0}
	//--------------------
	dst := image.NewPaletted(bounds, palette2)
	draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

	//images = append(images, dst)
	images[i] = dst
	//delays = append(delays, p.allCanvas[i].GifDelay)
	delays[i] = p.allCanvas[i].GifDelay
	//	disposals = append(disposals, gif.DisposalBackground)
	disposals[i] = gif.DisposalBackground

	mu.Lock()
	if img.Rect.Dx() > *gifWidth {
		*gifWidth = img.Rect.Dx()
	}
	if img.Rect.Dy() > *gifHeight {
		*gifHeight = img.Rect.Dy()
	}
	mu.Unlock()
}
