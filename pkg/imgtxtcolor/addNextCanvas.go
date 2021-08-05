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
	if p.canvas.img == nil {
		p.drw = &font.Drawer{
			Dst: p.canvas.img,
			Src: p.opt.FgColor,
			Face: truetype.NewFace(p.currentFont, &truetype.Options{
				Size:    float64(p.opt.FontSizeInt),
				Hinting: font.HintingFull,
			}),
		}
	}
	// создаем новый Canvas
	p.canvas = &ImgCanvas{
		padding: &stPadding{},
	}
	//	newC.img = p.canvas.img
	p.canvas.img = image.NewRGBA(image.Rect(0, 0, p.opt.Width, p.opt.Height)) // новая Canvas image.RGBA
	canvasSetBackground(p, color.Transparent)                                 // закрасим фон если надо добавим скругление углов
	p.drw.Dst = p.canvas.img                                                  // подключим ее
	p.textHeightSumm = fixed.I(0)                                             // сбросим курсор на 0 на первую строку сверху
	p.textWidthSumm = fixed.I(0)
	p.isNewCanvas = false // сбрасываем флаг, нам больше не требуется новый Canvas, его только что чоздали
	p.setFontSize(p.opt.FontSizeInt)
	// данные индивидуальные для этого Image
	p.canvas.bgColor = p.opt.BgColor
	p.canvas.alignVertical = p.opt.AlignVertical
	p.canvas.padding.bottom = p.padding.bottom
	p.canvas.padding.left = p.padding.left
	p.canvas.padding.top = p.padding.top
	p.canvas.padding.right = p.padding.right
	p.canvas.round = p.round
	p.canvas.gifDelay = p.opt.GifDelay
	// сразу сохраним Canvas в массиве

	p.allCanvas = append(p.allCanvas, p.canvas)
	// for i, u := range p.allCanvas {
	// 	log.Println("imgg", i, u.gifDelay)
	// }

}

// --------------
func drawCircle(p *stParam) {
	ctx := gg.NewContextForRGBA(p.canvas.img)
	ctx.SetColor(color.Opaque)
	ctx.DrawCircle(float64(p.padding.top), float64(p.padding.top/2), float64(p.padding.top/2)-2)
	ctx.Fill()
	//ctx.ResetClip()
}
