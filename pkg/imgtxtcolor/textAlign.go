package imgtxtcolor

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
)

func (p *stParam) textAlign() {
	p.textAlign2(p.canvas.Img)
}

func (p *stParam) textAlign2(img *image.RGBA) {
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
	// ctxTxt := gg.NewContext(b.Dx(), b.Dy())
	// ctxTxt.DrawImage(m, 0, 0)

	// вставим m text в новый Image
	draw.Draw(textCrop, textCrop.Bounds(), m, b.Min, draw.Src)
	// gradient
	if len(p.canvas.fontGradient) > 1 {
		ctxText := gg.NewContextForRGBA(textCrop)
		mask := ctxText.AsMask()
		//g := gg.NewLinearGradient(0, 0, float64(m.Bounds().Dx()), float64(m.Bounds().Dy()))
		gradient(ctxText, p.canvas.fontGradient, p.canvas.fontGradVector)
		//ctxText.SetFillStyle(g)
		// Используя маску, заливаем контекст градиентом
		ctxText.SetMask(mask)
		ctxText.DrawRectangle(0, 0, float64(m.Bounds().Dx()), float64(m.Bounds().Dy()))
		ctxText.Fill()
	}
	// end gradient

	canvasHeight := p.canvas.Img.Rect.Dy() // высота оригинала
	canvasWidth := p.canvas.Img.Rect.Dx()
	// вычисляем размеры
	if p.canvas.autoHeight {
		canvasHeight = p.canvas.padding.top + b.Dy() + p.canvas.padding.bottom
	}
	if p.canvas.autoWidth {
		canvasWidth = p.canvas.padding.left + (xRight - xLeft) + p.canvas.padding.right
	}

	// новый Image, если нужен
	if p.canvas.autoHeight || p.canvas.autoWidth {
		// минимальные ограничения
		if p.canvas.MinHeight > canvasHeight {
			canvasHeight = p.canvas.MinHeight
		}
		if p.canvas.MinWidth > canvasWidth {
			canvasWidth = p.canvas.MinWidth
		}
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
	if len(p.canvas.bgColor) > 1 {
		gradient(ctx, p.canvas.bgColor, p.canvas.bgGragVector)
	} else if len(p.canvas.bgColor) > 0 {
		ctx.SetColor(p.canvas.bgColor[0])
	}

	ctx.Fill()

	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(p.canvas.Img, p.canvas.Img.Bounds(), textCrop, pointSP, draw.Over)

	p.isNewCanvas = true // если будет новый текст создать в новом Image

}

// ищет начало пустоты в конце текста
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

// ищет начало текста сверху
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

// ищет начало текста слева
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

// ищет конец текста справа
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

// вырезает прямоугольник занятый текстом отбрасывая пустоту вокруг
//	на прозрачном фоне rgba = 0.0.0.0
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

func gradient(ctx *gg.Context, colors []color.RGBA, vector string) error {
	if len(colors) < 2 {
		return errors.New("в градиенте должно быть больше одного цвета")
	}
	step := float32(1) / float32(len(colors)-1)
	xTop := ctx.Image().Bounds().Dx() / 2
	var grad gg.Gradient
	if vector != "" {
		//parseVector(ctx, vector)
		grad = gg.NewLinearGradient(parseVector(ctx, vector))
	} else {
		grad = gg.NewLinearGradient(float64(xTop), 0, float64(xTop), float64(ctx.Image().Bounds().Dy()))
	}
	stepNum := float32(0)
	for i := 0; i < len(colors)-1; i++ {
		grad.AddColorStop(float64(stepNum), colors[i])
		stepNum += step
	}
	grad.AddColorStop(float64(1), colors[len(colors)-1])
	ctx.SetFillStyle(grad)
	return nil
}

//направление для градиента
//	разбор строки 0:0:0:0 в float64
func parseVector(ctx *gg.Context, vector string) (x0 float64, y0 float64, x1 float64, y1 float64) {
	a := strings.Split(vector, ":")
	rect := ctx.Image().Bounds()
	if len(a) == 4 {
		if f1, err := strconv.ParseFloat(a[0], 64); err == nil {
			if f2, err := strconv.ParseFloat(a[1], 64); err == nil {
				if f3, err := strconv.ParseFloat(a[2], 64); err == nil {
					if f4, err := strconv.ParseFloat(a[3], 64); err == nil {
						x0 := float64(rect.Dx()) / 100 * f1
						y0 := float64(rect.Dy()) / 100 * f2
						x1 := float64(rect.Dx()) / 100 * f3
						y1 := float64(rect.Dy()) / 100 * f4
						return x0, y0, x1, y1
					}
				}
			}
		}
	}
	log.Println("не удалось распарсить направление градиента")
	xTop := ctx.Image().Bounds().Dx() / 2
	return float64(xTop), 0, float64(xTop), float64(ctx.Image().Bounds().Dy())
}
