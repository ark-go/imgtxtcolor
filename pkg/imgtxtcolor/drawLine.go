package imgtxtcolor

import (
	"strings"

	"golang.org/x/image/math/fixed"
)

func (p *stParam) drawLine(text string) bool {

	if len(text) > 0 && (p.canvas == nil || p.isNewCanvas) {
		p.addNextCanvas()
	}

	p.drw.Src = p.opt.FgColor

	if strings.TrimSpace(text) == "" {
		//	log.Println("yes11")
		if ok := checkHeight(p, "WW", false); !ok {
			//		log.Println("yes")
			return true
		}
	}

	if ok := checkHeight(p, text, true); !ok { // у нас перебор по высоте
		if len(text) > 0 { // есть текст
			p.textToHeight()                           // выравнивание по вертикали напечатанный Image
			p.addNextCanvas()                          // добавим еще Canvas
			if ok := checkHeight(p, text, true); !ok { // и снова перебор, говорим об ошибке
				return false // нет места для строк
			}
		} else {
			// у нас перебор, но не влезает только пустая строка, нечего ее и печатать
			return true
		}
	}

	//log.Println("draw:", param.drw.Dot.Y)
	//log.Println("..", text)
	p.drw.DrawString(text)
	//p.isNewCanvas = false
	return true
}

// проверка влезает ли текст по вертикали,
//	в случае успеха передвигает указатель вертикальной печати если setHeight = true
//	setHeight = true - Y позиция будет изменена
//	setHeight = false - тест без реального передвижения указателя, только для проверки пустой строки
//	влезет ли что-нибудь после пустой строки, если не влезет, то и не будем рисовать пустую строку
//	или, тест на вместимость еще двух строк
func checkHeight(param *stParam, text string, setHeight bool) bool {
	param.drw.Dot.X = param.xPositionFunc(text) // выравнивание строки
	metric := param.drw.Face.Metrics()
	param.textHeightSumm += param.getHeight() // куда надо писать, отступ сверху (размер шрифта + lineSpacing)
	yPosition := fixed.I(param.padding.top)
	yPosition += param.textHeightSumm // куда надо писать, отступ сверху + paddingTop
	if !setHeight {
		// тест если текущая строка пустая проверим влезет ли следующая строка
		// yPosition - позиция Y куда писать текущую строку
		yPosition += param.getHeight() + metric.Descent // позиция после вставки текущей строки
		yPosition += param.getHeight() + metric.Descent // позиция для следущей, возможной строки
	}
	//! Добавил !!  param.getHeight()
	if yPosition.Ceil()+param.getHeight().Ceil()+metric.Descent.Ceil()+param.padding.bottom > param.canvas.Rect.Dy() {
		param.textHeightSumm -= param.getHeight() // если перебор, возвращаем позицию
		return false
	}

	param.drw.Dot.Y = yPosition // двигаем вниз

	return true
}
