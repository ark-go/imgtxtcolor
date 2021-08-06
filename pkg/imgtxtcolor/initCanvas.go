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

type alignVertical int
type alignHorizontal int

const (
	AlignVerticalTop alignVertical = iota + 1
	AlignVerticalCenter
	AlignVerticalBottom

	AlignHorizontalLeft alignHorizontal = iota + 1
	AlignHorizontalCenter
	AlignHorizontalRight
)

type ImgCanvas struct {
	img           *image.RGBA
	padding       *stPadding
	bgColor       color.RGBA
	round         float64
	alignVertical alignVertical
	gifDelay      int
}
type stParam struct {
	drw *font.Drawer
	// текуший Image
	//canvas    *image.RGBA
	canvas *ImgCanvas
	// Текущий шрифт
	currentFont *truetype.Font
	// Не реализовано // TODO
	border stBorder
	//
	textHeightSumm fixed.Int26_6 // для хранения высоты "курсора"
	textWidthSumm  fixed.Int26_6 // TODO
	textHeightTmp  fixed.Int26_6 // для расчета высоты строки
	// массив нарисованных Image
	allCanvas []*ImgCanvas //[]*image.RGBA
	palette   map[color.RGBA]bool
	// временно
	opt *stStartOptions
	// сигнал о том что следущий вывод требует нового Image
	isNewCanvas bool
}

func (p *stParam) setFontSize(size int) {
	if size < 1 {
		size = 20
	}
	p.opt.FontSize = size    // сохраним выбор
	if p.canvas.img != nil { // TODO первого может не быть? пока не появятся буквы мы не создаем Canvas
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
	FontSize int
	// цвет шрифта
	FgColor *image.Uniform
	// цвет фона
	BgColor color.RGBA
	// по высоте
	AlignVertical alignVertical
	// по горизонтали
	AlignHorizontal alignHorizontal
	// gif - file name если не задан не будет создаваться
	GifFileName string
	// gif delay  1/100 sec
	GifDelay int
	// paddings
	Padding *stPadding
	// Межстрочный интервал
	LineSpacing int
	// скругленные углы 0 - их нет
	Round float64
}

// Начальные установки по умолчанию
func StartOption() *stStartOptions {
	return &stStartOptions{
		Width:           500,
		Height:          350,
		FontSize:        20,
		FgColor:         &image.Uniform{C: colornames.Yellow},
		BgColor:         colornames.Darkslategray,
		AlignVertical:   AlignVerticalCenter,
		AlignHorizontal: AlignHorizontalCenter,
		GifFileName:     "",
		GifDelay:        100 * 4,
		Padding: &stPadding{
			left:   20,
			right:  20,
			top:    20,
			bottom: 20,
		},
		LineSpacing: 2,
		// скругленные углы
		Round: 0,
	}
}
func initCanvas(startOption *stStartOptions) (*stParam, error) {
	if startOption == nil {
		startOption = StartOption()
	}
	var err error
	param := stParam{
		canvas: &ImgCanvas{
			padding: &stPadding{},
		},
		border:         stBorder{false, 10, 10, 10, 10},
		textHeightSumm: fixed.I(0),
		textWidthSumm:  fixed.I(0),
		textHeightTmp:  fixed.I(0),
		allCanvas:      []*ImgCanvas{},
		palette:        make(map[color.RGBA]bool),
		opt:            startOption,
		isNewCanvas:    true,
	}

	setCurrentFont(&param, nil)

	param.canvas.img = nil
	param.drw = nil
	return &param, err
}

func canvasSetBackground(param *stParam, col color.Color) {
	ctx := gg.NewContextForRGBA(param.canvas.img)

	ctx.DrawRoundedRectangle(0, 0, float64(param.opt.Width), float64(param.opt.Height), float64(param.opt.Round))
	//ctx.SetColor(param.opt.BgColor)
	ctx.SetColor(col)
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
