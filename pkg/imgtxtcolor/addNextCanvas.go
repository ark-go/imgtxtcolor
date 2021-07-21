package imgtxtcolor

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func (p *stParam) addNextCanvas() {
	if p.canvas == nil {
		if fontFace, err := freetype.ParseFont(goregular.TTF); err == nil {
			p.drw = &font.Drawer{
				Dst: p.canvas,
				Src: p.startOption.fgColor,
				Face: truetype.NewFace(fontFace, &truetype.Options{
					Size:    float64(p.startOption.FontSizeInt),
					Hinting: font.HintingFull,
				}),
			}
		}
	}
	// создаем новый Canvas
	p.canvas = createCanvas(p) // новая Canvas image.RGBA
	canvasSetBackground(p) // закрасим фон если надо добавим скругление углов
	p.drw.Dst = p.canvas          // подключим ее
	p.textHeightSumm = fixed.I(0) // сбросим курсор на 0 на первую строку сверху
	p.textWidthSumm = fixed.I(0)
	// сразу сохраним Canvas в массиве
	p.allImages = append(p.allImages, p.canvas)
	p.isNewCanvas = false
	//addCircle(p)
	cl, _ := getColor("green")
	drawCircle(p, float64(p.padding.top/2), float64(p.padding.top/2), float64(p.padding.top/2), cl)
}

// --------------

func drawCircle(p *stParam, x, y, r float64, c color.RGBA) {
	ctx := gg.NewContextForRGBA(p.canvas)
	ctx.SetHexColor("ff00ff")
	ctx.DrawCircle(x, y, r)
	ctx.Fill()
	//ctx.SetHexColor("0000ff")
	// ctx.SetRGBA(0, 0, 0, 0)
	// ctx.DrawCircle(x, y, r/2)
	// ctx.Fill()
	//ctx := gg.NewContextForRGBA(p.canvas)

	ctx.DrawCircle(x, y*2, r)
	ctx.SetRGBA255(0, 255, 0, 100)
	//ctx.SetHexColor("ff00ff")

	ctx.Fill()
	//drawCircle2(p, x, y*2, r, c)
	// dc := gg.NewContext(1000, 1000)
	// dc.DrawCircle(350, 500, 300)
	// dc.Clip()
	// dc.DrawCircle(650, 500, 300)
	// dc.Clip()
	// // dc.DrawRectangle(0, 0, 1000, 1000)
	// dc.SetRGB(0, 0, 0)
	// dc.Fill()
	// dc.SavePNG("out.png")

	ctx.MoveTo(30, 30)
	ctx.LineTo(80, 30)
	ctx.LineTo(80, 150)
	//	ctx.DrawCircle(80, 150, 20)
	ctx.LineTo(30, 150)
	ctx.LineTo(30, 30)
	ctx.ClosePath()
	ctx.DrawCircle(80, 150, 20)
	ctx.SetRGBA(255, 0, 0, 100)
	//ctx.SetRGBA255(0, 255, 0, 100)
	ctx.Fill()

	ctx.DrawRoundedRectangle(50, 50, 300, 200, 20)
	ctx.SetRGBA(0, 250, 0, 100)
	ctx.Fill()
}
func drawCircle2(p *stParam, x, y, r float64, c color.RGBA) {
	ctx := gg.NewContextForRGBA(p.canvas)
	ctx.DrawCircle(x, y, r)
	ctx.SetRGBA255(100, 0, 255, 200)
	//ctx.SetHexColor("ff00ff")

	ctx.Fill()

}
