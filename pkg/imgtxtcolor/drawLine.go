package imgtxtcolor

import (
	"golang.org/x/image/math/fixed"
)

// отрисовываем строку
func (p *stParam) drawLine(text string) bool {

	if p.canvas.Img == nil || p.isNewCanvas {
		if text == "" {
			return true // новый Image и первое что пришло это пусто, может быть если были команды они подтирают за собой все пробелы
		}
		p.addNextCanvas()
	}

	p.drw.Src = p.opt.FgColor

	if ok := p.checkHeight(text); !ok { // у нас перебор по высоте
		if len(text) > 0 { // есть текст
			p.textAlign()                       // выравнивание по вертикали напечатанный Image
			p.addNextCanvas()                   // добавим еще Canvas
			if ok := p.checkHeight(text); !ok { // и снова перебор, говорим об ошибке
				return false // нет места для строк
			}
		} else {
			// у нас перебор, но не влезает только пустая строка, нечего ее и печатать
			return true
		}
	}
	p.canvas.setMaxX(p.drw.MeasureString(text))
	p.canvas.setmaxY(p.textHeightSumm) // TODO Нужно прибавлять Descent  + p.drw.Face.Metrics().Descent
	p.drw.DrawString(text)
	return true
}

// проверка влезает ли текст по вертикали,
//	в случае успеха передвигает указатель вертикальной печати
//	влезет ли что-нибудь после пустой строки, если не влезет, то и не будем рисовать пустую строку
//	или, тест на вместимость еще двух строк
func (p *stParam) checkHeight(text string) bool {
	p.drw.Dot.X = p.getHorizontalPos(text) // выравнивание строки
	metric := p.drw.Face.Metrics()
	// полная высота шрифта = межстрочный отступ / высота до базовой линии шрифта / крючки ниже линии шрифта
	fullHeightFont := fixed.I(p.opt.LineSpacing) + metric.Ascent + metric.Descent
	if p.textHeightSumm == 0 {
		// у первой строки нет Descent
		p.textHeightSumm += metric.Ascent // ascent - от верха до базовой линии шрифта
	} else {
		p.textHeightSumm += fullHeightFont // куда надо писать, отступ сверху (размер шрифта + lineSpacing)
	}
	yPosition := fixed.I(p.canvas.padding.top)
	yPosition += p.textHeightSumm // куда надо писать, отступ сверху + paddingTop

	if yPosition.Ceil()+p.canvas.padding.bottom > p.canvas.Img.Rect.Dy() {
		p.textHeightSumm -= fullHeightFont // если перебор, возвращаем позицию
		return false
	}

	p.drw.Dot.Y = yPosition // двигаем вниз

	return true
}
