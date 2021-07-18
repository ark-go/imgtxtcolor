package imgtxtcolor

import (
	"image"
	"image/draw"
	"log"
)

func textToBottomHeight(param *stParam) {
	if param.canvas == nil {
		log.Println("textToBottomHeight: Что-то не так.")
		return
	}
	textHeight := param.textHeightSumm.Ceil() + param.padding.top + param.drw.Face.Metrics().Descent.Ceil() // + param.drw.Face.Metrics().Descent.Ceil() // текущая линия,                   //
	m := param.canvas.SubImage(image.Rect(0, param.padding.top, param.width, textHeight))                   //.(*image.RGBA)

	// перевод его в новый RGB  чтоб закрасить старый и наложить обратно   // TODO почему image.image не так рисуется как image.RGBA
	// размеры отрезанного куска
	b := m.Bounds()
	newCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// или вставляем наш отрезанный
	draw.Draw(newCrop, newCrop.Bounds(), m, b.Min, draw.Src)
	// color, _ := getColor("yellow")
	// draw.Draw(newCrop, newCrop.Bounds(), &image.Uniform{C: color}, image.Point{}, draw.Src)
	iHeight := param.canvas.Rect.Dy()                    // высота Canvas
	iHeightTxt := textHeight - param.padding.top         // нижняя граница текста минус top padding
	top := (iHeight - iHeightTxt) - param.padding.bottom // все свободное место минус padding-bottom
	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := param.canvas.Bounds().Min.Add(image.Point{0, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в newCrop
	draw.Draw(param.canvas, param.canvas.Bounds(), &image.Uniform{C: param.bgColor}, image.Point{}, draw.Src)
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(param.canvas, param.canvas.Bounds(), newCrop, pointSP, draw.Src)
	//!param.textHeightSumm = fixed.I(param.height) // если дальше будет текст он сам  создаст новый Image
	//addNextCanvas(param) // новый канвас, не нужен тут , т.е. в конце всего текста появится пустой Image
	//param.textHeightSumm = fixed.I(0)
	param.isNewCanvas = true

}
