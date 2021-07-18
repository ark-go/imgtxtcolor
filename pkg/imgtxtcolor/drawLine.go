package imgtxtcolor

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

func drawLine(param *stParam, text string) bool {

	if param.canvas == nil || param.isNewCanvas {
		if len(text) > 0 || param.canvas == nil {
			param.addNextCanvas()
		}
	}

	fontFace, _ := freetype.ParseFont(goregular.TTF)
	param.drw.Face = truetype.NewFace(fontFace, &truetype.Options{
		Size:    float64(param.startOption.FontSizeInt),
		Hinting: font.HintingFull,
	})

	param.drw.Src = param.startOption.fgColor
	if ok := checkHeight(param, text); !ok {
		if len(text) > 0 {
			param.addNextCanvas()
			if ok := checkHeight(param, text); !ok {
				return false // нет места для строк
			}
			//			log.Println("height text", param.textHeightSumm.Ceil(), text)
		} else {
			// у нас перебор, но не влезает пустая строка, нечего печатать
			return true
		}
	}

	//log.Println("draw:", param.drw.Dot.Y)
	param.drw.DrawString(text)
	return true
}
func checkHeight(param *stParam, text string) bool {
	param.drw.Dot.X = param.xPositionFunc(text)
	metric := param.drw.Face.Metrics()
	param.textHeightSumm += param.getHeight()
	yPosition := fixed.I(param.padding.top)
	yPosition += param.textHeightSumm
	if yPosition.Ceil()+metric.Descent.Ceil()+param.padding.bottom > param.canvas.Rect.Dy() {
		return false
	}
	param.drw.Dot.Y = yPosition // двигаем вниз
	return true
}
