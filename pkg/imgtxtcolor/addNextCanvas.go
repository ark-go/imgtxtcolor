package imgtxtcolor

import (
	"image"
	//	"image/color"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (p *stParam) addNextCanvas() {

	//	if a < 1 {
	// создаем новый Canvas
	p.canvas = &ImgCanvas{
		parent:  p,
		padding: &stPadding{},
	}
	//p.canvas.Img = image.NewRGBA(image.Rect(0, 0, 5, 5))

	p.canvas.Img = image.NewRGBA(image.Rect(0, 0, p.opt.Width, p.opt.Height)) // новая Canvas image.RGBA

	p.drw = &font.Drawer{
		Dst: p.canvas.Img, // подключим ее
		Src: &image.Uniform{C: p.opt.FontColor[0]},
		Face: truetype.NewFace(p.currentFont, &truetype.Options{
			Size:    float64(p.opt.FontSize),
			Hinting: font.HintingFull,
		}),
	}
	//	}
	//	a = a + 1
	//canvasSetBackground(p, color.Transparent) // TODO: (здесь уже не нужно наверно) закрасим фон если надо добавим скругление углов
	// инициализация данных Image
	p.textHeightSumm = fixed.I(0) // сбросим курсор Y на 0 на первую строку сверху
	//	p.setFontSize(p.opt.FontSize)
	// данные индивидуальные для этого Image
	p.canvas.bgColor = p.opt.BgColor
	p.canvas.alignVertical = p.opt.AlignVertical
	p.canvas.alignHorizontal = p.opt.AlignHorizontal
	p.canvas.padding.bottom = p.opt.Padding.bottom
	p.canvas.padding.left = p.opt.Padding.left
	p.canvas.padding.top = p.opt.Padding.top
	p.canvas.padding.right = p.opt.Padding.right
	p.canvas.round = p.opt.Round
	p.canvas.GifDelay = p.opt.GifDelay
	p.canvas.autoHeight = p.opt.AutoHeight
	p.canvas.autoWidth = p.opt.AutoWidth
	p.canvas.MinWidth = p.opt.MinWidth
	p.canvas.MinHeight = p.opt.MinHeight
	p.canvas.fontColor = p.opt.FontColor
	p.canvas.fontGradVector = p.opt.FontGradVector
	p.canvas.bgGragVector = p.opt.BgGragVector
	p.canvas.frameFilePath = p.opt.FrameFilePath
	// сразу сохраним Canvas в массиве
	p.allCanvas = append(p.allCanvas, p.canvas)
	// сбрасываем флаг, нам больше не требуется новый Canvas, его только что чоздали
	p.isNewCanvas = false
	//log.Printf("Time [%v]: %v\n", "Image создан за:", time.Since(startTime))

}

// --------------
// func drawCircle(p *stParam) {
// 	ctx := gg.NewContextForRGBA(p.canvas.Img)
// 	ctx.SetColor(color.Opaque)
// 	ctx.DrawCircle(float64(p.opt.Padding.top), float64(p.opt.Padding.top/2), float64(p.opt.Padding.top/2)-2)
// 	ctx.Fill()
// 	//ctx.ResetClip()
// }
