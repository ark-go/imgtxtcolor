package imgtxtcolor

import (
	"image"
	"image/draw"

	"github.com/fogleman/gg"
)

func (p *stParam) textAlignVertical() {
	if p.canvas != nil && p.canvas.img == nil {
		// если приходим сюда без Canvas, значит это первый вход
		// и еще не было текста, а только пустые строки, поэтому Canvas еще не создавался
		// пустая строка это пустая, не пробел, если в строке были команды и за ними перевод строки
		log.Println("textAlignVertical: Image еще нет, не забыть отключить это сообщение.")
		return
	}

	testY := p.getBottomBorder()
	textHeight := testY //
	x0 := p.canvas.img.Rect.Min.X + p.canvas.padding.left
	y0 := p.canvas.img.Rect.Min.Y + p.canvas.padding.top
	x1 := p.canvas.img.Rect.Max.X - p.canvas.padding.right
	y2 := textHeight
	m := p.canvas.img.SubImage(image.Rect(x0, y0, x1, y2))
	// в m не новое изображение, там теже пиксели, а в Bounds размеры которые мы устанвливали
	b := m.Bounds()
	// создадим новый Image, в него будем переносить наш m
	newCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// вставим m в новый Image
	draw.Draw(newCrop, newCrop.Bounds(), m, b.Min, draw.Src)
	var top int
	switch p.canvas.alignVertical {
	case AlignVerticalCenter:
		iHeight := p.canvas.img.Rect.Dy()               // высота оригинала
		iHeightTxt := textHeight - p.canvas.padding.top // высота текста без верхнего padding
		//iHeightTxt := testY
		top = (iHeight - iHeightTxt) / 2 // половина свободного места
	case AlignVerticalBottom:
		iHeight := p.canvas.img.Rect.Dy()                      // высота Canvas
		iHeightTxt := textHeight - p.canvas.padding.top        // нижняя граница текста минус top padding
		top = (iHeight - iHeightTxt) - p.canvas.padding.bottom // все свободное место минус padding-bottom
	case AlignVerticalTop:
		top = p.canvas.padding.top
	}
	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := p.canvas.img.Bounds().Min.Add(image.Point{p.canvas.padding.left * -1, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в newCrop
	ctx := gg.NewContextForRGBA(p.canvas.img)
	ctx.Clear()
	ctx.DrawRoundedRectangle(0, 0, float64(p.canvas.img.Rect.Dx()), float64(p.canvas.img.Rect.Dy()), float64(p.canvas.round))
	ctx.SetColor(p.canvas.bgColor)
	ctx.Fill()
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(p.canvas.img, p.canvas.img.Bounds(), newCrop, pointSP, draw.Over)

	p.isNewCanvas = true // если будет новый текст создать в новом Image

}
func (p *stParam) getBottomBorder() int {
	// At(Bounds().Min.X, Bounds().Min.Y)   возвращает верхний левый пиксель сетки.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1)  возвращает нижний правый.
	//var u = color.RGBA{0, 0, 0, 0} // Transparent
	var h int
	for i := p.canvas.img.Bounds().Max.Y - 1; i > p.canvas.img.Bounds().Min.Y; i-- { // Y
		for j := p.canvas.img.Bounds().Min.X; j < p.canvas.img.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := p.canvas.img.At(j, i).RGBA(); r == 0 && g == 0 && b == 0 && a == 0 {
				h++
			} else {
				return i + 1
			}
		}
	}
	// log.Println("пустые:", h, h2)
	return 0
}
