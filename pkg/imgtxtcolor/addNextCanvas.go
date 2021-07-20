package imgtxtcolor

import (
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
	p.canvas = createCanvas(p)    // новая Canvas image.RGBA
	p.drw.Dst = p.canvas          // подключим ее
	p.textHeightSumm = fixed.I(0) // сбросим курсор на 0 на первую строку сверху
	p.textWidthSumm = fixed.I(0)
	// сразу сохраним Canvas в массиве
	p.allImages = append(p.allImages, p.canvas)
	p.isNewCanvas = false
}
