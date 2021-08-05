package imgtxtcolor

import (
	"image"
	"image/draw"

	"github.com/fogleman/gg"
)

func textToCenterHeight(param *stParam) {
	if param.canvas == nil {
		log.Println("textToCenterHeight: Что-то не так.")
		return
	}
	testY := param.getBottomBorder()
	//textHeight := param.textHeightSumm.Ceil() + param.padding.top + param.drw.Face.Metrics().Descent.Ceil() // текущая линия,
	textHeight := testY //
	x0 := param.canvas.Rect.Min.X + param.canvasOpt.padding.left
	y0 := param.canvas.Rect.Min.Y + param.canvasOpt.padding.top
	x1 := param.canvas.Rect.Max.X - param.canvasOpt.padding.right
	y2 := textHeight
	m := param.canvas.SubImage(image.Rect(x0, y0, x1, y2))
	// в m не новое изображение, там теже пиксели, а в Bounds размеры которые мы устанвливали
	b := m.Bounds()
	// создадим новый Image, в него будем переносить наш m
	newCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// вставим m в новый Image
	draw.Draw(newCrop, newCrop.Bounds(), m, b.Min, draw.Src)

	iHeight := param.canvas.Rect.Dy()                      // высота оригинала
	iHeightTxt := textHeight - param.canvasOpt.padding.top // высота текста без верхнего padding
	//iHeightTxt := testY
	top := (iHeight - iHeightTxt) / 2 // половина свободного места
	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := param.canvas.Bounds().Min.Add(image.Point{param.canvasOpt.padding.left * -1, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в newCrop
	ctx := gg.NewContextForRGBA(param.canvas)
	ctx.Clear()
	ctx.DrawRoundedRectangle(0, 0, float64(param.canvas.Rect.Dx()), float64(param.canvas.Rect.Dy()), float64(param.round))
	ctx.SetColor(param.canvasOpt.bgColor)
	ctx.Fill()
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(param.canvas, param.canvas.Bounds(), newCrop, pointSP, draw.Over)

	param.isNewCanvas = true // если будет новый текст создать в новом Image

}
func (p *stParam) getBottomBorder() int {
	// At(Bounds().Min.X, Bounds().Min.Y)   возвращает верхний левый пиксель сетки.
	// At(Bounds().Max.X-1, Bounds().Max.Y-1)  возвращает нижний правый.

	//var u = color.RGBA{0, 0, 0, 0}
	var h int
	for i := p.canvas.Bounds().Max.Y - 1; i > p.canvas.Bounds().Min.Y; i-- { // Y
		for j := p.canvas.Bounds().Min.X; j < p.canvas.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := p.canvas.At(j, i).RGBA(); r == 0 && g == 0 && b == 0 && a == 0 {
				h++
			} else {
				return i + 1
			}
		}
	}
	// log.Println("пустые:", h, h2)
	return 0
}
