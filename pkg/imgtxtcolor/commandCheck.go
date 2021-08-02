package imgtxtcolor

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"

	"golang.org/x/image/math/fixed"
)

func commandCheck(param *stParam, str, cmd string) bool {
	isCmd := true
	// str = strings.ToLower(str)  color,,???
	switch strings.ToLower(cmd) {
	case "fontsize", "size":
		if siz, err := strconv.Atoi(str); err == nil {
			param.setFontSize(siz)

		} else {
			isCmd = false
		}
		// fallthrough // Переходит на следующий иначе break
	case "fontcolor", "color":
		//if col, ok := colornames.Map[str]; ok {
		if col, ok := getColor(str); ok {
			//param.drw.Src = &image.Uniform{C: col}
			param.opt.FgColor = &image.Uniform{C: col}
			param.palette[col] = true
			// if param.canvas != nil {
			// 	param.drw.Src = param.opt.FgColor
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
		default:
			isCmd = false // не засчитали команду
		}

	case "padding":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.setAll(siz)
		} else {
			isCmd = false
		}
	case "paddingtop", "top":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.top = siz
		} else {
			isCmd = false
		}
	case "paddingleft", "left":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.left = siz
		} else {
			isCmd = false
		}
	case "paddingrigh", "right":
		if siz, err := strconv.Atoi(str); err == nil {
			param.padding.right = siz
		} else {
			isCmd = false
		}
	case "paddingbottom", "bottom":
		if siz, err := strconv.Atoi(str); err == nil && siz >= 0 {
			param.padding.bottom = siz
		} else {
			isCmd = false
		}
	case "round":
		if siz, err := strconv.Atoi(str); err == nil && siz >= 0 {
			param.round = float64(siz)
		} else {
			if param.padding.top > 0 {
				param.round = float64(param.padding.top)
			}
		}
	case "linespacing":
		if siz, err := strconv.Atoi(str); err == nil {
			param.lineSpacing = fixed.I(siz)
		} else {
			isCmd = false
		}
	case "bgcolor":
		// Только в начале текста, иначе все закрасит
		if col, ok := getColor(str); ok {
			param.textToHeight()
			param.opt.BgColor = col
			param.palette[col] = true
		} else {
			isCmd = false
		}
	case "width":
		if siz, err := strconv.Atoi(str); err == nil {
			param.opt.Width = siz
			param.isNewCanvas = true
		} else {
			isCmd = false
		}
	case "height":
		if siz, err := strconv.Atoi(str); err == nil {
			// замена Canvas
			if len(param.allImages) > 0 {
				param.allImages = param.allImages[:len(param.allImages)-1]
			}
			param.opt.Height = siz
			param.isNewCanvas = true
		} else {
			isCmd = false
		}
	case "rect":
		if str == "tg" {
			width, height := getRectToTelegram(float64(param.opt.Width), float64(param.opt.Height))
			param.opt.Width, param.opt.Height = int(width), int(height)
			log.Println(width, height)

		} else {
			isCmd = false
		}
	case "break", "alignh":
		switch strings.ToLower(str) {
		case "top":
			param.opt.AlignHeight = "top"
		case "center":
			param.opt.AlignHeight = "center"
			param.textToHeight()
		case "bottom":
			param.opt.AlignHeight = "bottom"
			param.textToHeight()
		default:
			isCmd = false // не засчитали команду
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
