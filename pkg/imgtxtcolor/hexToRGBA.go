package imgtxtcolor

import (
	"fmt"
	"image/color"
)

// смотреть еще  https://github.com/pwaller/go-hexcolor/blob/master/hexcolor.go

func hexToRGBA(hex string) (color.RGBA, error) {
	var (
		rgba             color.RGBA
		err              error
		errInvalidFormat = fmt.Errorf("invalid")
	)
	rgba.A = 0xff
	if hex[0] != '#' {
		return rgba, errInvalidFormat
	}
	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}
	switch len(hex) {
	case 7:
		rgba.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		rgba.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		rgba.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
	case 4:
		rgba.R = hexToByte(hex[1]) * 17
		rgba.G = hexToByte(hex[2]) * 17
		rgba.B = hexToByte(hex[3]) * 17
	case 5:
		rgba.R = hexToByte(hex[1]) * 17
		rgba.G = hexToByte(hex[2]) * 17
		rgba.B = hexToByte(hex[3]) * 17
		rgba.A = hexToByte(hex[4]) * 17
	case 9:
		rgba.R = hexToByte(hex[1])<<4 + hexToByte(hex[2])
		rgba.G = hexToByte(hex[3])<<4 + hexToByte(hex[4])
		rgba.B = hexToByte(hex[5])<<4 + hexToByte(hex[6])
		rgba.A = hexToByte(hex[7])<<4 + hexToByte(hex[8])
	default:
		err = errInvalidFormat

	}
	return rgba, err
}
