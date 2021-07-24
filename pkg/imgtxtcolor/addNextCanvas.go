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
	if p.canvas == nil {
		p.drw = &font.Drawer{
			Dst: p.canvas,
			Src: p.opt.FgColor,
			Face: truetype.NewFace(p.currentFont, &truetype.Options{
				Size:    float64(p.opt.FontSizeInt),
				Hinting: font.HintingFull,
			}),
		}
	}
	// создаем новый Canvas
	p.canvas = image.NewRGBA(image.Rect(0, 0, p.opt.Width, p.opt.Height)) // новая Canvas image.RGBA
	canvasSetBackground(p)                                                // закрасим фон если надо добавим скругление углов
	p.drw.Dst = p.canvas                                                  // подключим ее
	p.textHeightSumm = fixed.I(0)                                         // сбросим курсор на 0 на первую строку сверху
	p.textWidthSumm = fixed.I(0)
	// сразу сохраним Canvas в массиве
	p.allImages = append(p.allImages, p.canvas)
	p.isNewCanvas = false // сбрасываем флаг, нам больше не требуется новый Canvas, его только что чоздали
	p.setFontSize(p.opt.FontSizeInt)
	//addCircle(p)
	//	drawCircle(p)
}

// --------------
func drawCircle(p *stParam) {
	ctx := gg.NewContextForRGBA(p.canvas)
	ctx.SetColor(color.Opaque)
	ctx.DrawCircle(float64(p.padding.top), float64(p.padding.top/2), float64(p.padding.top/2)-2)
	ctx.Fill()
	//ctx.ResetClip()
}
