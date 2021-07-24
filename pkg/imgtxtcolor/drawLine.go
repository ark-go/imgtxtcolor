package imgtxtcolor

import (
	"golang.org/x/image/math/fixed"
)

func (p *stParam) drawLine(text string) bool {

	if p.canvas == nil || p.isNewCanvas { // nil такой ситуации не должно быть
		if len(text) > 0 || p.canvas == nil {
			p.addNextCanvas()
		}
	}

	p.drw.Src = p.opt.FgColor
	if ok := checkHeight(p, text); !ok { // у нас перебор по высоте
		if len(text) > 0 {
			p.addNextCanvas()                    // добавим еще Canvas
			if ok := checkHeight(p, text); !ok { // и снова перебор, говорим об ошибке
				return false // нет места для строк
			}
		} else {
			// у нас перебор, но не влезает только пустая строка, нечего ее и печатать
			return true
		}
	}

	//log.Println("draw:", param.drw.Dot.Y)
	p.drw.DrawString(text)
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
