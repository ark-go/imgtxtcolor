package imgtxtcolor

import (
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
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
	parent *stParam
	// Canvas
	Img *image.RGBA
	// padding
	padding *stPadding
	// цвет фона
	bgColor        []color.RGBA
	bgGragVector   string
	fontColor      []color.RGBA
	fontGradVector string
	// радиус углов
	round float64
	// Вертикальное выравнивание
	alignVertical alignVertical
	// Горизонтальное выравнивание
	alignHorizontal alignHorizontal
	// задержка кадра, если использовать в gif / webp
	GifDelay int
	// максимальная ширина строки в этом Canvas
	maxX fixed.Int26_6
	// максимальная высота текста в этом Canvas
	maxY fixed.Int26_6
	// автоматическая ширина по тексту
	autoWidth bool
	// автоматическая высота по тексту
	autoHeight bool
	// минимальная ширина
	MinWidth int
	// минимальная высота
	MinHeight int
	// frame path  путь к рамке
	frameFilePath string
}

// если val больше AllMaxX заменяем
func (c *ImgCanvas) setMaxX(val fixed.Int26_6) {
	if val > c.maxX {
		c.maxX = val
	}
	if c.maxX > c.parent.allMaxX {
		c.parent.allMaxX = c.maxX
	}
}

// если val больше AllMaxY заменяем
func (c *ImgCanvas) setmaxY(val fixed.Int26_6) {
	if val > c.maxY {
		c.maxY = val
	}
	if c.maxY > c.parent.allMaxY {
		c.parent.allMaxY = c.maxY
	}
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
	textHeightTmp  fixed.Int26_6 // для расчета высоты строки
	// массив нарисованных Image
	allCanvas []*ImgCanvas //[]*image.RGBA
	palette   map[color.RGBA]bool
	// временно
	opt *stStartOptions
	// сигнал о том что следущий вывод требует нового Image
	isNewCanvas bool
	allMaxX     fixed.Int26_6
	allMaxY     fixed.Int26_6
}

// установка шрифта
func (p *stParam) setFontSize(size int) {
	if size < 1 {
		size = 20
	}
	p.opt.FontSize = size    // сохраним выбор
	if p.canvas.Img != nil { // TODO первого может не быть? пока не появятся буквы мы не создаем Canvas
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
	FontColor []color.RGBA
	// цвет фона
	BgColor        []color.RGBA
	BgGragVector   string
	FontGradVector string
	// по высоте
	AlignVertical alignVertical
	// по горизонтали
	AlignHorizontal alignHorizontal
	// gif file name path если не задан? gif не будет создаваться
	GifFileName string
	// gif delay  1/100 sec
	GifDelay int
	// paddings
	Padding *stPadding
	// Межстрочный интервал
	LineSpacing int
	// скругленные углы 0 - их нет
	Round float64
	// авто Height, ширина по тексту
	AutoHeight bool
	// авто Width, высота по тексту
	AutoWidth bool
	// минимальная ширина
	MinWidth int
	// минимальная высота
	MinHeight int
	// Каталог для фреймов рамок
	FrameDir string
	// Файл рамки, текущий, убедитесь что путь существует!!!
	FrameFilePath string
}

// Начальные установки по умолчанию
func StartOption() *stStartOptions {
	rootDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln("не определить рабочий каталог")

	}

	return &stStartOptions{
		Width:           500,
		Height:          350,
		MinWidth:        5,
		MinHeight:       5,
		FontSize:        20,
		FontColor:       []color.RGBA{colornames.Yellow},
		BgColor:         []color.RGBA{colornames.Darkslategray},
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
		//путь
		FrameDir:      filepath.Join(rootDir, "frame"),
		FrameFilePath: "",
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
		textHeightTmp:  fixed.I(0),
		allCanvas:      []*ImgCanvas{},
		palette:        make(map[color.RGBA]bool),
		opt:            startOption,
		isNewCanvas:    true,
	}

	param.setCurrentFont(nil)

	param.canvas.Img = nil
	param.drw = nil
	return &param, err
}

func canvasSetBackground(param *stParam, col color.Color) {
	ctx := gg.NewContextForRGBA(param.canvas.Img)
	ctx.DrawRoundedRectangle(0, 0, float64(param.opt.Width), float64(param.opt.Height), float64(param.opt.Round))
	ctx.SetColor(col)
	ctx.Fill()
}

// Установка Font (goregular.TTF)
// 	font установить в nil  // TODO сделать возможность выбора шрифта
func (p *stParam) setCurrentFont(font []byte) {
	if font == nil {
		font = goregular.TTF
		_ = gobold.TTF
	}
	if fontFace, err := freetype.ParseFont(font); err == nil {
		p.currentFont = fontFace
	} else {
		log.Println("Ошибка при загрузке шрифта")
		p.currentFont = nil
	}
}
