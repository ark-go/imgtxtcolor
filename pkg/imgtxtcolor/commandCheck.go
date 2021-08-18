package imgtxtcolor

import (
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

// на входе значение и команда
//	Выход:
//	_cmd - подтверждает что это существующая команда
//	_break - команда требует создания нового Image
func (p *stParam) commandCheck(str, cmd string) (_cmd, _break bool) {
	isCmd := true
	// str = strings.ToLower(str)  color,,???
	if strings.HasPrefix(strings.ToLower(cmd), "padding") {
		siz, err := strconv.Atoi(str)
		if err != nil || siz < 0 {
			return false, false
		}
		switch strings.ToLower(cmd) {
		case "padding":
			p.opt.Padding.setAll(siz)
		case "paddingtop":
			p.opt.Padding.top = siz
		case "paddingleft":
			p.opt.Padding.left = siz
		case "paddingright":
			p.opt.Padding.right = siz
		case "paddingbottom":
			p.opt.Padding.bottom = siz
		default:
			return false, false
		}
		return true, true
	}
	//----------------------------------------------
	switch strings.ToLower(cmd) {
	case "fontsize", "size":
		if siz, err := strconv.Atoi(str); err == nil {
			p.setFontSize(siz)
			return true, false
		}
		// fallthrough // Переходит на следующий иначе break
	case "fontcolor", "color":
		// массив для градиента, если единственный элемент - цвет шрифта
		if col, v, ok := getColorArr(p, str); ok {
			p.opt.FontColor = col    // это массив
			p.opt.FontGradient = col // и это массив
			p.opt.FontGradVector = v
			return true, false
		}
		if str == "transparent" {
			p.opt.FontColor = []color.RGBA{{0, 0, 0, 0}}
			return true, false
		}
	case "align":
		// определяем функции для расчета позиции по горизонтали
		switch strings.ToLower(str) {
		case "left":
			p.opt.AlignHorizontal = AlignHorizontalLeft
		case "center":
			p.opt.AlignHorizontal = AlignHorizontalCenter
		case "right":
			p.opt.AlignHorizontal = AlignHorizontalRight
		default:
			isCmd = false // не засчитали команду
		}
		return isCmd, false
	case "round":
		if siz, err := strconv.ParseFloat(str, 64); err == nil && siz >= 0 {
			p.opt.Round = float64(siz)
		} else {
			if p.opt.Padding.top > 0 {
				p.opt.Round = float64(p.opt.Padding.top)
			}
		}
		return true, true
	case "delay":
		if siz, err := strconv.Atoi(str); err == nil && siz >= 0 {
			p.opt.GifDelay = siz
			return true, true
		}
		return false, false
	case "linespacing":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.LineSpacing = siz
			return true, false
		}
	case "bgcolor":
		// Только в начале текста, иначе все закрасит
		// массив для градиента, если единственный элемент - цвет фона
		if col, v, ok := getColorArr(p, str); ok {
			p.isNewCanvas = true
			p.opt.BgColor = col // это массив
			p.opt.BgGragVector = v
			//p.palette[col] = true
			return true, true
		}
		if str == "transparent" {
			p.opt.BgColor = []color.RGBA{{0, 0, 0, 0}}
			return true, true
		}
		return false, false
	case "width":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Width = siz
			p.opt.AutoWidth = false
			return true, true
		}
	case "height":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Height = siz
			p.opt.AutoHeight = false
			return true, true
		}
	case "maxwidth":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Width = siz
			p.opt.AutoWidth = true
			return true, true
		}
	case "maxheight":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Height = siz
			p.opt.AutoHeight = true
			return true, true
		}
	case "minwidth":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.MinWidth = siz
			return true, true
		}
	case "minheight":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.MinHeight = siz
			return true, true
		}
	case "rect":
		if str == "tg" {
			width, height := getRectToTelegram(float64(p.opt.Width), float64(p.opt.Height))
			p.opt.Width, p.opt.Height = int(width), int(height)
			log.Println("rect-tg", width, height)
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
			p.opt.AlignVertical = AlignVerticalTop
			return true, true
		case "center":
			p.opt.AlignVertical = AlignVerticalCenter
			return true, true
		case "bottom":
			p.opt.AlignVertical = AlignVerticalBottom
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

func getColorArr(p *stParam, str string) ([]color.RGBA, string, bool) {
	colArr := strings.Split(str, ":")

	colArrRgba := []color.RGBA{}
	if len(colArr) > 0 {
		for i, val := range colArr {
			if val == "" { // два двоеточия рядом породят пустую строку
				return []color.RGBA{}, "", false
			}
			// проверим на символ X он разрешает писать команды к градиенту после цветов
			if []rune(val)[0] == 'V' { // vector

				if len(colArrRgba) > 1 {
					vec := colArr[i:]
					if len(vec) == 5 {
						return colArrRgba, strings.Join(vec[1:], ":"), true
					} else {
						return []color.RGBA{}, "", false
					}
				}
			}
			// добавляем цвет если смогли его определить
			if colR, ok := getColor(val); ok {
				colArrRgba = append(colArrRgba, colR)
				continue
			} else {
				return []color.RGBA{}, "", false
			}
		}
		return colArrRgba, "", true
	} else {
		return []color.RGBA{}, "", false
	}

}

// Приводим название цвета или его код # ffffff  к color.RGBA
//	если не получится то возвращаем false
func getColor(str string) (color.RGBA, bool) {
	if str == "" {
		return color.RGBA{}, false
	}
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
