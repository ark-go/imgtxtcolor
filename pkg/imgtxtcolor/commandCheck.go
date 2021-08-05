package imgtxtcolor

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"

	"golang.org/x/image/math/fixed"
)

func commandCheck(param *stParam, str, cmd string) (_cmd, _break bool) {
	isCmd := true
	// str = strings.ToLower(str)  color,,???
	if strings.HasPrefix(strings.ToLower(cmd), "padding") {
		siz, err := strconv.Atoi(str)
		if err != nil && siz < 0 {
			return false, false
		}
		switch strings.ToLower(cmd) {
		case "padding":
			param.padding.setAll(siz)
		case "paddingtop":
			param.padding.top = siz
		case "paddingleft":
			param.padding.left = siz
		case "paddingright":
			param.padding.right = siz
		case "paddingbottom":
			param.padding.bottom = siz
		default:
			return false, false
		}
		return true, true
	}
	//----------------------------------------------
	switch strings.ToLower(cmd) {
	case "fontsize", "size":
		if siz, err := strconv.Atoi(str); err == nil {
			param.setFontSize(siz)
			return true, false
		}
		// fallthrough // Переходит на следующий иначе break
	case "fontcolor", "color":
		//if col, ok := colornames.Map[str]; ok {
		if col, ok := getColor(str); ok {
			param.opt.FgColor = &image.Uniform{C: col}
			param.palette[col] = true
			return true, false
		}
	case "align":
		// определяем функции для расчета позиции по горизонтали
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
		return isCmd, false
	case "round":
		if siz, err := strconv.Atoi(str); err == nil && siz >= 0 {
			param.round = float64(siz)
		} else {
			if param.padding.top > 0 {
				param.round = float64(param.padding.top)
			}
		}
		return true, true
	case "linespacing":
		if siz, err := strconv.Atoi(str); err == nil {
			param.lineSpacing = fixed.I(siz)
			return true, false
		}
	case "bgcolor":
		// Только в начале текста, иначе все закрасит
		if col, ok := getColor(str); ok {
			param.isNewCanvas = true
			param.opt.BgColor = col
			param.palette[col] = true
			return true, true
		}
		return false, false
	case "width":
		if siz, err := strconv.Atoi(str); err == nil {
			param.opt.Width = siz
			return true, true
		}
	case "height":
		if siz, err := strconv.Atoi(str); err == nil {
			// замена Canvas
			// if len(param.allImages) > 0 {
			// 	param.allImages = param.allImages[:len(param.allImages)-1]
			// }
			param.opt.Height = siz
			return true, true
		}
	case "rect":
		if str == "tg" {
			width, height := getRectToTelegram(float64(param.opt.Width), float64(param.opt.Height))
			param.opt.Width, param.opt.Height = int(width), int(height)
			log.Println(width, height)
			return true, true
		}
	case "break":
		switch strings.ToLower(str) {
		case "page":
			return true, true
		default:
			isCmd = false // не засчитали команду
		}
		return isCmd, false
	case "alignv", "alignh":
		switch strings.ToLower(str) {
		case "top":
			param.opt.AlignVertical = AlignVerticalTop
			return true, true
		case "center":
			param.opt.AlignVertical = AlignVerticalCenter
			return true, true
		case "bottom":
			param.opt.AlignVertical = AlignVerticalBottom
			return true, true
		default:
			isCmd = false // не засчитали команду
		}
		return isCmd, false
	default:
		isCmd = false // не засчитали команду

	}
	return isCmd, false
}

func getColor(str string) (color.RGBA, bool) {
	if str[0] == '#' {
		if col, err := hexToRGBA(str); err != nil {
			return color.RGBA{}, false
		} else {
			return col, true
		}
	}
	if str == "transparent" {
		var u = color.RGBA{0, 0, 0, 0}
		return u, true
	}
	col1, ok := colornames.Map[str]
	return col1, ok

}
