package imgtxtcolor

import (
	"image"
	"image/draw"

	"github.com/fogleman/gg"
)

func textToTopHeight(param *stParam) {
	if param.canvas == nil {
		log.Println("textToBottomHeight: Что-то не так.")
		return
	}

	textHeight := param.textHeightSumm.Ceil() + param.padding.top + param.drw.Face.Metrics().Descent.Ceil() // + param.drw.Face.Metrics().Descent.Ceil() // текущая линия,                   //

	x0 := param.canvas.Rect.Min.X + param.padding.left
	y0 := param.canvas.Rect.Min.Y + param.padding.top
	x1 := param.canvas.Rect.Max.X - param.padding.right
	//y2 := param.canvas.Rect.Max.Y
	y2 := textHeight
	m := param.canvas.SubImage(image.Rect(x0, y0, x1, y2))

	//	m := param.canvas.SubImage(image.Rect(param.padding.left, param.padding.top, param.opt.Width-param.padding.right, textHeight))

	// перевод его в новый RGB  чтоб закрасить старый и наложить обратно
	// размеры отрезанного куска
	b := m.Bounds()
	newCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// или вставляем наш отрезанный
	draw.Draw(newCrop, newCrop.Bounds(), m, b.Min, draw.Src)
	//iHeight := param.canvas.Rect.Dy()                    // высота Canvas
	//iHeightTxt := textHeight - param.padding.top         // нижняя граница текста минус top padding
	//top := (iHeight - iHeightTxt) - param.padding.bottom // все свободное место минус padding-bottom
	top := param.canvasOpt.padding.top
	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := param.canvas.Bounds().Min.Add(image.Point{param.canvasOpt.padding.left * -1, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в newCrop
	ctx := gg.NewContextForRGBA(param.canvas)
	ctx.DrawRoundedRectangle(0, 0, float64(param.canvas.Rect.Dx()), float64(param.canvas.Rect.Dy()), float64(param.round))
	ctx.SetColor(param.canvasOpt.bgColor)
	ctx.Fill()
	//draw.Draw(param.canvas, param.canvas.Bounds(), &image.Uniform{C: param.opt.BgColor}, image.Point{}, draw.Src)
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(param.canvas, param.canvas.Bounds(), newCrop, pointSP, draw.Over) //Over Src

	param.isNewCanvas = true // сообщаем - требуется новый Canvas

}
