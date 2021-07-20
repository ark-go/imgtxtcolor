package imgtxtcolor

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"

	"golang.org/x/image/math/fixed"
)

func commandCheck(param *stParam, str, cmd string) bool {
	isCmd := true
	// str = strings.ToLower(str)  color,,???
	switch cmd {
	case "fontSize", "size":
		if siz, err := strconv.Atoi(str); err == nil {
			param.startOption.FontSizeInt = siz
			if param.canvas != nil {
				fontFace, _ := freetype.ParseFont(goregular.TTF)
				param.drw.Face = truetype.NewFace(fontFace, &truetype.Options{
					Size:    float64(siz),
					Hinting: font.HintingFull,
				})
			}
		}
		// fallthrough // Переходит на следующий иначе break
	case "fontColor", "color":
		//if col, ok := colornames.Map[str]; ok {
		if col, ok := getColor(str); ok {
			//param.drw.Src = &image.Uniform{C: col}
			param.startOption.fgColor = &image.Uniform{C: col}
			// if param.canvas != nil {
			// 	param.drw.Src = param.startOption.fgColor
			// }
		} else {
			isCmd = false
		}
	case "align":
		switch strings.ToLower(str) {
		case "left":
			param.xPositionFunc = func(str string) fixed.Int26_6 { return fixed.I(param.padding.left) } // влево
		case "center":
			param.xPositionFunc = func(str string) fixed.Int26_6 {
				max := fixed.I(param.canvas.Rect.Max.X) // всего
				max -= fixed.I(param.padding.right)     // отнимаем справа
				max -= fixed.I(param.padding.left)      // отнимаем слева
				max -= param.drw.MeasureString(str)     // получаем свободное место
				max /= 2                                //место пополам
				max += fixed.I(param.padding.left)      // отодвигаем слева
				return max
				//	return (fixed.I(param.canvas.Rect.Max.X) + fixed.Int26_6(param.padding.left) - param.drw.MeasureString(str)) / 2
			} // для центра
		case "right":
			param.xPositionFunc = func(str string) fixed.Int26_6 {
				return (fixed.I(param.canvas.Rect.Dx()) - param.drw.MeasureString(str) - fixed.I(param.padding.right))
			} // вправо
		}

	case "padding":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.setAll(siz)
		}
	case "paddingTop", "top":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.top = siz
		}
	case "paddingLeft", "left":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.left = siz
		}
	case "paddingRigh", "right":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.right = siz
		}
	case "paddingBottom", "bottom":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.bottom = siz
		}
	case "lineSpacing", "linespacing":
		if siz, err := strconv.Atoi(str); err == nil {
			param.lineSpacing = fixed.I(siz)
		}
	case "bgColor", "bgcolor":
		// Только в начале текста, иначе все закрасит
		if col, ok := getColor(str); ok {
			param.bgColor = col
			// if param.canvas != nil {
			// 	draw.Draw(param.canvas, param.canvas.Bounds(), &image.Uniform{C: col},
			// 		image.Point{}, draw.Src)
			// }
		} else {
			isCmd = false
		}
	case "width":
		if siz, err := strconv.Atoi(str); err == nil {
			// замена Canvas
			// if len(param.allImages) > 0 {
			// 	param.allImages = param.allImages[:len(param.allImages)-1]
			// }
			param.width = siz
			param.isNewCanvas = true
			//param.addNextCanvas() // новая канва
		}
	case "height":
		if siz, err := strconv.Atoi(str); err == nil {
			// замена Canvas
			if len(param.allImages) > 0 {
				param.allImages = param.allImages[:len(param.allImages)-1]
			}
			param.height = siz
			param.isNewCanvas = true
			//param.addNextCanvas() // новая ка
		}

	case "break", "position":
		switch strings.ToLower(str) {
		case "top":
		case "center":
			textToCenterHeight(param)
		case "bottom":
			textToBottomHeight(param)
		}
	default:
		isCmd = false // не засчитали команду

	}
	return isCmd
}

func getColor(str string) (color.RGBA, bool) {
	if str[0] == '#' {
		if col, err := hexToRGBA(str); err != nil {
			return color.RGBA{}, false
		} else {
			return col, true
		}
	}
	col1, ok := colornames.Map[str]
	return col1, ok

}
