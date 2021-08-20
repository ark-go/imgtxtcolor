package imgtxtcolor

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"strconv"
	"strings"
	"sync"

	"github.com/fogleman/gg"
)

// func (p *stParam) textAlign() {

// 	p.isNewCanvas = true // если будет новый текст создать в новом Image
// }

// вызов горутин для окончательного форматирования всех созданных Image
// до вызова все Image (p.allCanvas) содержат текст на прозрачном фоне, в размерах maxWidth/maxHeight или width/height
func (p *stParam) formatAllCanvas() {
	var wg sync.WaitGroup
	for _, canvas := range p.allCanvas {
		wg.Add(1)
		go canvas.formatCanvas(&wg)
	}
	wg.Wait()
}

// окончательное форматирование Image
//	до вызова Image содержит текст на прозрачном фоне, в размерах maxWidth/maxHeight или width/height
//	установка фона, градиента, позиции текста в соответсвии с командами в тексте
func (canvas *ImgCanvas) formatCanvas(wg *sync.WaitGroup) {
	defer wg.Done()
	if canvas != nil && canvas.Img == nil {
		// если приходим сюда без Canvas, значит это первый вход
		// и еще не было текста, а только пустые строки, поэтому Canvas еще не создавался
		// пустая строка это пустая, не пробел, если в строке были команды и за ними перевод строки
		return
	}

	yTop, yBottom, xLeft, xRight := getRectText(canvas.Img)
	textHeight := yBottom - yTop
	// получаем координаты нашего текста, вырезаем текст
	m := canvas.Img.SubImage(image.Rect(xLeft, yTop, xRight, yBottom))
	// в m не новое изображение, там те же пиксели, а в Bounds размеры которые мы устанвливали
	b := m.Bounds()

	// создадим новый Image, для нашего вырезанного текста
	textCrop := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	// ctxTxt := gg.NewContext(b.Dx(), b.Dy())
	// ctxTxt.DrawImage(m, 0, 0)

	// вставим m text в новый Image
	draw.Draw(textCrop, textCrop.Bounds(), m, b.Min, draw.Src)
	// gradient
	if len(canvas.fontColor) > 1 {
		ctxText := gg.NewContextForRGBA(textCrop)
		mask := ctxText.AsMask()
		//g := gg.NewLinearGradient(0, 0, float64(m.Bounds().Dx()), float64(m.Bounds().Dy()))
		gradient(ctxText, canvas.fontColor, canvas.fontGradVector)
		//ctxText.SetFillStyle(g)
		// Используя маску, заливаем контекст градиентом
		ctxText.SetMask(mask)
		ctxText.DrawRectangle(0, 0, float64(m.Bounds().Dx()), float64(m.Bounds().Dy()))
		ctxText.Fill()
	}
	// end gradient

	canvasHeight := canvas.Img.Rect.Dy() // высота оригинала
	canvasWidth := canvas.Img.Rect.Dx()
	// вычисляем размеры
	if canvas.autoHeight {
		canvasHeight = canvas.padding.top + b.Dy() + canvas.padding.bottom
	}
	if canvas.autoWidth {
		canvasWidth = canvas.padding.left + (xRight - xLeft) + canvas.padding.right
	}

	// новый Image, если нужен
	if canvas.autoHeight || canvas.autoWidth {
		// минимальные ограничения
		if canvas.MinHeight > canvasHeight {
			canvasHeight = canvas.MinHeight
		}
		if canvas.MinWidth > canvasWidth {
			canvasWidth = canvas.MinWidth
		}
		// меняем размер оригинала под новый размер текста
		canvas.Img = image.NewRGBA(image.Rect(0, 0, canvasWidth, canvasHeight))
	}

	var top int
	switch canvas.alignVertical {
	case AlignVerticalCenter:
		//iHeight := canvasHeight          // высота оригинала
		iHeight := canvasHeight - canvas.padding.top - canvas.padding.bottom // рабочая зона
		//	iHeightTxt := textHeight         // - canvas.padding.top // высота текста
		top = (iHeight - textHeight) / 2 // половина свободного места TODO:
		top += canvas.padding.top
	case AlignVerticalBottom:
		//	iHeight := canvasHeight                              // высота Canvas
		//	iHeightTxt := textHeight                             // - canvas.padding.top        //
		top = canvasHeight - textHeight - canvas.padding.bottom // все место - размер текста - padding-bottom
	case AlignVerticalTop:
		top = canvas.padding.top
	}
	var offsetLeft int
	switch canvas.alignHorizontal {
	case AlignHorizontalCenter:
		iWidth := canvasWidth - canvas.padding.left - canvas.padding.right
		iWidthTxt := textCrop.Rect.Dx()
		offsetLeft = (iWidth - iWidthTxt) / 2
		offsetLeft += canvas.padding.left
	//offsetLeft = (canvas.Img.Rect.Dx() - textCrop.Rect.Dx()) / 2 // TODO:
	case AlignHorizontalRight:
		offsetLeft = canvas.Img.Rect.Dx() - textCrop.Rect.Dx() - canvas.padding.right
	case AlignHorizontalLeft:
		offsetLeft = canvas.padding.left
	}

	// точка для совмещения нашего отрезанного куска с основным изображением
	pointSP := canvas.Img.Bounds().Min.Add(image.Point{offsetLeft * -1, top * -1})
	// теперь закрасим все, а все что надо мы уже отрезали в textCrop
	ctx := gg.NewContextForRGBA(canvas.Img)
	ctx.Clear()
	ctx.DrawRoundedRectangle(0, 0, float64(canvasWidth), float64(canvasHeight), float64(canvas.round))
	if len(canvas.bgColor) > 1 {
		gradient(ctx, canvas.bgColor, canvas.bgGragVector)
	} else if len(canvas.bgColor) > 0 {
		ctx.SetColor(canvas.bgColor[0])
	}

	ctx.Fill()

	// /home/arkadii/ProjectsGo/canvas/frame408x292.png
	if canvas.frameFilePath != "" {
		setFrame(canvas, canvas.frameFilePath) //frame408x292.png
		//setFrameAroundImg(canvas, canvas.frameFilePath)
	}
	// sp? 4-й параметр. точка совмещения она совместится с dest в точке 0.0 и все съедет относительно точки
	draw.Draw(canvas.Img, canvas.Img.Bounds(), textCrop, pointSP, draw.Over)
}

// ищет начало пустоты в конце текста
func getBottomText(img *image.RGBA, c chan int) {
	for i := img.Bounds().Max.Y - 1; i > img.Bounds().Min.Y; i-- { // Y
		for j := img.Bounds().Min.X; j < img.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := img.At(j, i).RGBA(); (r + g + b + a) != 0 {
				c <- i + 1
				return
			}
		}
	}
	c <- 0
}

// ищет начало текста сверху
func getTopText(img *image.RGBA, c chan int) {
	for i := img.Bounds().Min.Y; i < img.Bounds().Max.Y-1; i++ { // Y
		for j := img.Bounds().Min.X; j < img.Bounds().Max.X-1; j++ { // X
			if r, g, b, a := img.At(j, i).RGBA(); (r + g + b + a) != 0 {
				c <- i
				return
			}
		}
	}
	c <- 0
}

// ищет начало текста слева
func getLeftText(img *image.RGBA, c chan int) {
	for i := img.Bounds().Min.X; i < img.Bounds().Max.X-1; i++ { // X
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y-1; j++ {
			if r, g, b, a := img.At(i, j).RGBA(); (r + g + b + a) != 0 {
				c <- i
				return
			}
		}
	}
	c <- 0
}

// ищет конец текста справа
func getRightText(img *image.RGBA, c chan int) {
	for i := img.Bounds().Max.X - 1; i > img.Bounds().Min.X; i-- { // X
		for j := img.Bounds().Min.Y; j < img.Bounds().Max.Y-1; j++ {
			if r, g, b, a := img.At(i, j).RGBA(); (r + g + b + a) != 0 {
				c <- i + 1
				return
			}
		}
	}
	c <- 0
}

// вырезает прямоугольник занятый текстом отбрасывая пустоту вокруг
//	на прозрачном фоне rgba = 0.0.0.0
func getRectText(img *image.RGBA) (yTop, yBottom, xLeft, xRight int) {

	top := make(chan int)
	bottom := make(chan int)
	left := make(chan int)
	right := make(chan int)
	go getTopText(img, top)
	go getBottomText(img, bottom)
	go getLeftText(img, left)
	go getRightText(img, right)

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
