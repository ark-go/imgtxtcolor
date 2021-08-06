package imgtxtcolor

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func (p *stParam) addNextCanvas() {
	// создаем новый Canvas
	p.canvas = &ImgCanvas{
		padding: &stPadding{},
	}
	p.canvas.img = image.NewRGBA(image.Rect(0, 0, p.opt.Width, p.opt.Height)) // новая Canvas image.RGBA
	p.drw = &font.Drawer{
		Dst: p.canvas.img, // подключим ее
		Src: p.opt.FgColor,
		Face: truetype.NewFace(p.currentFont, &truetype.Options{
			Size:    float64(p.opt.FontSize),
			Hinting: font.HintingFull,
		}),
	}
	canvasSetBackground(p, color.Transparent) // закрасим фон если надо добавим скругление углов
	// инициализация данных Image
	p.textHeightSumm = fixed.I(0) // сбросим курсор Y на 0 на первую строку сверху
	p.textWidthSumm = fixed.I(0)  // TODO
	//	p.setFontSize(p.opt.FontSize)
	// данные индивидуальные для этого Image
	p.canvas.bgColor = p.opt.BgColor
	p.canvas.alignVertical = p.opt.AlignVertical
	p.canvas.padding.bottom = p.opt.Padding.bottom
	p.canvas.padding.left = p.opt.Padding.left
	p.canvas.padding.top = p.opt.Padding.top
	p.canvas.padding.right = p.opt.Padding.right
	p.canvas.round = p.opt.Round
	p.canvas.gifDelay = p.opt.GifDelay
	// сразу сохраним Canvas в массиве
	p.allCanvas = append(p.allCanvas, p.canvas)
	// сбрасываем флаг, нам больше не требуется новый Canvas, его только что чоздали
	p.isNewCanvas = false
}

// --------------
func drawCircle(p *stParam) {
	ctx := gg.NewContextForRGBA(p.canvas.img)
	ctx.SetColor(color.Opaque)
	ctx.DrawCircle(float64(p.opt.Padding.top), float64(p.opt.Padding.top/2), float64(p.opt.Padding.top/2)-2)
	ctx.Fill()
	//ctx.ResetClip()
}
