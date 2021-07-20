package imgtxtcolor

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type stBorder struct {
	isBorderOn bool
	top        int
	bottom     int
	left       int
	right      int
}

type stPadding struct {
	top    int
	bottom int
	left   int
	right  int
}

func (p *stPadding) lenW() int {
	return p.left + p.right
}
func (p *stPadding) setAll(v int) {
	p.left = v
	p.right = v
	p.top = v
	p.bottom = v
}

type stParam struct {
	drw *font.Drawer
	// текуший Image
	canvas *image.RGBA
	// структура top,bottom,left,right
	padding *stPadding
	// ширина Image
	width int
	// Высота Image
	height int
	// Не реализовано // TODO
	border stBorder
	// динамическая текущая функция для выравнивания текста, изменяется в зависимости от параметров
	xPositionFunc func(str string) fixed.Int26_6
	// базовая линия для рисования текста, Descent рисует ниже
	textHeightSumm fixed.Int26_6 // для хранения высоты "курсора"
	textWidthSumm  fixed.Int26_6 // TODO
	textHeightTmp  fixed.Int26_6 // для расчета высоты строки
	// межстрочное расстояние default: 2
	lineSpacing fixed.Int26_6
	// Цвет background
	bgColor color.RGBA
	// массив нарисованных Image
	allImages []*image.RGBA
	// временно
	startOption *stStartOptions
	// сигнал о том что следущий вывод требует нового Image
	isNewCanvas bool
}

// высота полного шрифта + межстрочный
func (p *stParam) getHeight() fixed.Int26_6 {
	if p.drw == nil {
		log.Println("error getHeight")
		return fixed.I(0)
	}
	metric := p.drw.Face.Metrics()
	//log.Println("height:", metric.Height, "AsDes", metric.Ascent+metric.Descent)
	return metric.Height + p.lineSpacing
	//return metric.Ascent + metric.Descent + p.lineSpacing // fixed.I(2) // рекомендуемое  metric.Height
}

// Начальные установки по умолчанию
type stStartOptions struct {
	// Ширина
	Width int
	// Высота
	Height int
	// Размер шрифта
	FontSizeInt int
	fgColor     *image.Uniform
}

// Начальные установки по умолчанию
func StartOption() *stStartOptions {
	return &stStartOptions{
		Width:       500,
		Height:      350,
		FontSizeInt: 20,
		fgColor:     &image.Uniform{C: colornames.Yellow},
	}
}
func initCanvas(opt *stStartOptions) (*stParam, error) {
	//	bgColor := colornames.Darkslategray
	//	fgColor := &image.Uniform{C: colornames.Yellow}
	//	fontSize := float64(opt.FontSizeInt)
	padding := stPadding{20, 20, 20, 20}
	//	var fontFace *truetype.Font
	var err error
	param := stParam{
		xPositionFunc:  func(str string) fixed.Int26_6 { return fixed.I(padding.left) }, // влево,
		padding:        &padding,
		border:         stBorder{false, 10, 10, 10, 10},
		textHeightSumm: fixed.I(0),
		textWidthSumm:  fixed.I(0),
		textHeightTmp:  fixed.I(0),
		lineSpacing:    fixed.I(2),
		bgColor:        colornames.Darkslategray,
		allImages:      []*image.RGBA{},
		width:          opt.Width,
		height:         opt.Height,
		startOption:    opt,
		isNewCanvas:    true,
	}
	//	fontFace, err = freetype.ParseFont(goregular.TTF)

	// param.canvas = createCanvas(&param)
	// param.drw = &font.Drawer{
	// 	Dst: param.canvas,
	// 	Src: opt.fgColor,
	// 	Face: truetype.NewFace(fontFace, &truetype.Options{
	// 		Size:    fontSize,
	// 		Hinting: font.HintingFull,
	// 	}),
	// }
	param.canvas = nil
	param.drw = nil
	return &param, err
}
func createCanvas(param *stParam) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, param.width, param.height))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: param.bgColor},
		image.Point{}, draw.Src)
	//param.drw.Dst = param.canvas
	// newR := image.Rect(30, 30, width-30, height-30)
	// draw.Draw(canvas, newR, &image.Uniform{C: colornames.Red},
	// 	image.Point{}, draw.Src)
	return canvas
}

/*
for i, r := range []rune(s) {
    fmt.Printf("%d: %q\n", i, r)
}
*/
