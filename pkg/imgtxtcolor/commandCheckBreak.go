package imgtxtcolor

import (
	"strconv"
	"strings"
)

func commandCheckBreak(param *stParam, str, cmd string) bool {
	isBreak := true
	str = strings.ToLower(str)
	switch cmd {
	case "break":
		switch strings.ToLower(str) {
		case "top":
		case "center":
			textToCenterHeight(param)
		case "bottom":
			textToBottomHeight(param)
		default:
			isBreak = false
		}
	case "break-width":
		if siz, err := strconv.Atoi(str); err == nil {
			param.width = siz
			//	param.height = siz // TODO !!
		} else {
			isBreak = false
		}

	// fallthrough // Переходит на следующий иначе break
	default:
		isBreak = false
	}
	return isBreak
}
