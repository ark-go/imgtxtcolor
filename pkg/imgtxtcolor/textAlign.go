package imgtxtcolor

import (
	"image"
	"image/draw"

	"github.com/fogleman/gg"
)

func (p *stParam) textAlign() {
	if p.canvas != nil && p.canvas.Img == nil {
		// если приходим сюда без Canvas, значит это первый вход
		// и еще не было текста, а только пустые строки, поэтому Canvas еще не создавался
		// пустая строка это пустая, не пробел, если в строке были команды и за ними перевод строки
		return
	}

	yTop, yBottom, xLeft, xRight := getRectText(p)
	textHeight := yBottom - yTop
	// получаем координаты нашего текста, вырезаем текст
	m := p.canvas.Img.SubImage(image.Rect(xLeft, yTop, xRight, yBottom))
	// в m не новое изображение, там те же пиксели, а в Bounds размеры которые мы устанвливали
	b := m.Bounds()

	// создадим новый Image, для нашего вырезанного текста
	textCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// вставим m в новый Image
	draw.Draw(textCrop, textCrop.Bounds(), m, b.Min, draw.Src)

	canvasHeight := p.canvas.Img.Rect.Dy() // высота оригинала
	canvasWidth := p.canvas.Img.Rect.Dx()
	if p.canvas.autoHeight {
		canvasHeight = p.canvas.padding.top + b.Dy() + p.canvas.padding.bottom

	}
	if p.canvas.autoWidth {
		canvasWidth = p.canvas.padding.left + (xRight - xLeft) + p.canvas.padding.right
	}
	if p.canvas.autoHeight || p.canvas.autoWidth {
		// меняем размер оригинала под новый размер текста
		p.canvas.Img = image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	}

	var top int
	switch p.canvas.alignVertical {
	case AlignVerticalCenter:
		iHeight := canvasHeight          // высота оригинала
		iHeightTxt := textHeight         // - p.canvas.padding.top // высота текста без верхнего padding
		top = (iHeight - iHeightTxt) / 2 // половина свободного места
	case AlignVerticalBottom:
		iHeight := canvasHeight                                // высота Canvas
		iHeightTxt := textHeight                               // - p.canvas.padding.top        // нижняя граница текста минус top padding
		top = (iHeight - iHeightTxt) - p.canvas.padding.bottom // все свободное место минус padding-bottom
	case AlignVerticalTop:
		top = p.canvas.padding.top
	}
	var offsetLeft int
	switch p.canvas.alignHorizontal {
	case AlignHorizontalCenter:
		offsetLeft = (p.canvas.Img.Rect.Dx() - textCrop.Rect.Dx()) / 2
	case AlignHorizontalRight:
		offsetLeft = p.canvas.Img.Rect.Dx() - textCrop.Rect.Dx() - p.canvas.padding.right
	case AlignHorizontalLeft:
		offsetLeft = p.canvas.padding.left
	}

	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := p.canvas.Img.Bounds().Min.Add(image.Point{offsetLeft * -1, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в textCrop
	ctx := gg.NewContextForRGBA(p.canvas.Img)
	ctx.Clear()
	ctx.DrawRoundedRectangle(0, 0, float64(canvasWidth), float64(canvasHeight), float64(p.canvas.round))
	ctx.SetColor(p.canvas.bgColor)
	ctx.Fill()
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(p.canvas.Img, p.canvas.Img.Bounds(), textCrop, pointSP, draw.Over)

	p.isNewCanvas = true // если будет новый текст создать в новом Image

}

// поиск не пустых ячеек снизу, т.е. окончание напечатанного текста
// func (p *stParam) getBottomBorder() int {
// 	// At(Bounds().Min.X, Bounds().Min.Y)   возвращает верхний левый пиксель сетки.
// 	// At(Bounds().Max.X-1, Bounds().Max.Y-1)  возвращает нижний правый.
// 	//var u = color.RGBA{0, 0, 0, 0} // Transparent
// 	//if p.canvas.autoHeight {
// 	// впринципе этот блок можно убрать совсем
// 	// здесь мы дополнительно удалим лишние пустые строки
// 	// но мы не //TODO
// 	var h int
// 	for i := p.canvas.Img.Bounds().Max.Y - 1; i > p.canvas.Img.Bounds().Min.Y; i-- { // Y
// 		for j := p.canvas.Img.Bounds().Min.X; j < p.canvas.Img.Bounds().Max.X-1; j++ { // X
// 			if r, g, b, a := p.canvas.Img.At(j, i).RGBA(); r == 0 && g == 0 && b == 0 && a == 0 {
// 				h++
// 			} else {
// 				//log.Println("y:", i+1)
// 				return i + 1
// 			}
// 		}
// 	}
// 	// log.Println("пустые:", h, h2)
// 	return 0
// 	// } else {
// 	// 	 здесь по размер по размеру строк, вместе с пустыми, расчитано по мере печати
// 	// 	y := p.canvas.maxY + p.drw.Face.Metrics().Descent
// 	// 	log.Println("y:", y.Ceil())
// 	// 	return y.Ceil() + p.canvas.padding.top
// 	// }
// }

func getBottomText(p *stParam, c chan int) {
	for i := p.canvas.Img.Bounds().Max.Y - 1; i > p.canvas.Img.Bounds().Min.Y; i-- { // Y
		for j := p.canvas.Img.Bounds().Min.X; j < p.canvas.Img.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := p.canvas.Img.At(j, i).RGBA(); (r + g + b + a) != 0 {
				c <- i + 1
				return
			}
		}
	}
	c <- 0
}
func getTopText(p *stParam, c chan int) {
	for i := p.canvas.Img.Bounds().Min.Y; i < p.canvas.Img.Bounds().Max.Y-1; i++ { // Y
		for j := p.canvas.Img.Bounds().Min.X; j < p.canvas.Img.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := p.canvas.Img.At(j, i).RGBA(); (r + g + b + a) != 0 {
				c <- i
				return
			}
		}
	}
	c <- 0
}
func getLeftText(p *stParam, c chan int) {
	for i := p.canvas.Img.Bounds().Min.X; i < p.canvas.Img.Bounds().Max.X-1; i++ { // X
		for j := p.canvas.Img.Bounds().Min.Y; j < p.canvas.Img.Bounds().Max.Y-1; j++ {
			if r, g, b, a := p.canvas.Img.At(i, j).RGBA(); (r + g + b + a) != 0 {
				c <- i
				return
			}
		}
	}
	c <- 0
}

func getRightText(p *stParam, c chan int) {
	for i := p.canvas.Img.Bounds().Max.X - 1; i > p.canvas.Img.Bounds().Min.X; i-- { // X
		for j := p.canvas.Img.Bounds().Min.Y; j < p.canvas.Img.Bounds().Max.Y-1; j++ {
			if r, g, b, a := p.canvas.Img.At(i, j).RGBA(); (r + g + b + a) != 0 {
				c <- i + 1
				return
			}
		}
	}
	c <- 0
}

func getRectText(p *stParam) (yTop, yBottom, xLeft, xRight int) {

	top := make(chan int)
	bottom := make(chan int)
	left := make(chan int)
	right := make(chan int)
	go getTopText(p, top)
	go getBottomText(p, bottom)
	go getLeftText(p, left)
	go getRightText(p, right)

	for i := 0; i < 4; i++ {
		// ждем два канала
		select {
		case yT := <-top:
			yTop = yT
		case yB := <-bottom:
			yBottom = yB
		case xL := <-left:
			xLeft = xL
		case xR := <-right:
			xRight = xR
		}
	}
	return
}

// func getNewDy(p *stParam){

// }
