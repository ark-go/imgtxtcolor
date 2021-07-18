package imgtxtcolor

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
)

func GetBase64(img *image.RGBA) string {
	// Буфер в памяти для хранения изображений PNG
	// прежде, чем мы закодируем базу 64
	var buff bytes.Buffer
	// Буфер соответствует интерфейсу Writer, поэтому мы можем использовать его с Encode
	png.Encode(&buff, img)
	// Закодируйте байты в буфере в строку base64
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())
	// Вы можете встроить его в html-документ с помощью этой строки
	//	htmlImage := "<img src=\"data:image/png;base64," + encodedString + "\" />"
	return encodedString
}
