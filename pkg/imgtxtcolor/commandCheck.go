package imgtxtcolor

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	"golang.org/x/image/colornames"
)

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
		if col, ok := getColor(str); ok {
			p.opt.FgColor = &image.Uniform{C: col}
			p.palette[col] = true
			return true, false
		}
		if str == "transparent" {
			p.opt.FgColor = &image.Uniform{C: color.RGBA{0, 0, 0, 0}}
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
		if col, ok := getColor(str); ok {
			p.isNewCanvas = true
			p.opt.BgColor = col
			p.palette[col] = true
			return true, true
		}
		if str == "transparent" {
			p.opt.BgColor = color.RGBA{0, 0, 0, 0}
			return true, true
		}
		return false, false
	case "width":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Width = siz
			return true, true
		}
	case "height":
		if siz, err := strconv.Atoi(str); err == nil {
			p.opt.Height = siz
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
