package internal

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"sync"

	"github.com/ark-go/imgtxtcolor/pkg/imgtxtcolor"
)

func getBase64(imgArr []*imgtxtcolor.ImgCanvas) []string {
	tmp := make([]string, len(imgArr))
	var wg sync.WaitGroup
	for i, canvas := range imgArr {
		wg.Add(1)
		go addBase64(tmp, i, canvas.Img, &wg)
	}
	wg.Wait()
	
	return tmp
}

func addBase64(tmp []string, i int, img *image.RGBA, wg *sync.WaitGroup) {
	defer wg.Done()
	var buff bytes.Buffer
	png.Encode(&buff, img)
	tmp[i] = base64.StdEncoding.EncodeToString(buff.Bytes())
}
