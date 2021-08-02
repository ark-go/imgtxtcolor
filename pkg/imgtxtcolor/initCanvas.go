package imgtxtcolor

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
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
	padding     *stPadding
	currentFont *truetype.Font
	// скругленные углы
	round float64
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
	// массив нарисованных Image
	allImages []*image.RGBA
	palette   map[color.RGBA]bool
	// временно
	opt *stStartOptions
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

func (p *stParam) setFontSize(size int) {
	if size < 1 {
		size = 20
	}
	p.opt.FontSizeInt = size // сохраним выбор
	if p.canvas != nil {     // TODO первого может не быть? пока не появятся буквы мы не создаем Canvas
		p.drw.Face = truetype.NewFace(p.currentFont, &truetype.Options{
			Size:    float64(size),
			Hinting: font.HintingFull,
		})
	}
}

// Начальные установки по умолчанию
type stStartOptions struct {
	// Ширина
	Width int
	// Высота
	Height int
	// Размер шрифта
	FontSizeInt int
	// цвет шрифта
	FgColor *image.Uniform
	// цвет фона
	BgColor color.RGBA
	// по высоте
	AlignHeight string
	// gif - file name
	GifFileName string
	// gif delay
	GifDelay int // 1/100 sec
}

// Начальные установки по умолчанию
func StartOption() *stStartOptions {
	return &stStartOptions{
		Width:       500,
		Height:      350,
		FontSizeInt: 20,
		FgColor:     &image.Uniform{C: colornames.Yellow},
		BgColor:     colornames.Darkslategray,
		AlignHeight: "center",
		GifFileName: "",
		GifDelay:    100 * 4,
	}
}
func initCanvas(startOption *stStartOptions) (*stParam, error) {
	if startOption == nil {
		startOption = StartOption()
	}
	padding := stPadding{20, 20, 20, 20}
	var err error
	param := stParam{
		xPositionFunc:  func(str string) fixed.Int26_6 { return fixed.I(padding.left) }, // влево,
		padding:        &padding,
		border:         stBorder{false, 10, 10, 10, 10},
		textHeightSumm: fixed.I(0),
		textWidthSumm:  fixed.I(0),
		textHeightTmp:  fixed.I(0),
		lineSpacing:    fixed.I(2),
		//	bgColor:        colornames.Darkslategray,
		allImages: []*image.RGBA{},
		palette:   make(map[color.RGBA]bool),
		//width:          startOption.Width,
		//height:      startOption.Height,
		opt:         startOption,
		isNewCanvas: true,
	}

	setCurrentFont(&param, nil)

	param.canvas = nil
	param.drw = nil
	return &param, err
}

func canvasSetBackground(param *stParam) {
	ctx := gg.NewContextForRGBA(param.canvas)

	ctx.DrawRoundedRectangle(0, 0, float64(param.opt.Width), float64(param.opt.Height), float64(param.round))
	ctx.SetColor(param.opt.BgColor)
	ctx.Fill()
}

// Установка Font (goregular.TTF)
// 	font установить в nil  // TODO сделать возможность выбора шрифта
func setCurrentFont(param *stParam, font []byte) {
	if font == nil {
		font = goregular.TTF
	}
	if fontFace, err := freetype.ParseFont(font); err == nil {
		param.currentFont = fontFace
	} else {
		log.Println("Ошибка при загрузке шрифта")
		param.currentFont = nil
	}
}
