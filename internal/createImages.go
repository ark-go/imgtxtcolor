package internal

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/ark-go/imgtxtcolor/pkg/imgtxtcolor"
)

func createImages(text string, gif bool) ([]*imgtxtcolor.ImgCanvas, error) {
	dir, _ := ioutil.ReadDir(imgDir)
	for _, d := range dir {
		os.RemoveAll(path.Join(imgDir, d.Name()))
	}
	opt := imgtxtcolor.StartOption()
	// opt.Width = sizeW
	// opt.Height = sizeH
	// opt.FontSize = fontSizeInt
	if gif {
		opt.GifFileName = path.Join(imgDir, "test.gif")
	}
	opt.GifDelay = 100 * 1

	//canvasArr, err := imgtxtcolor.CreateImageText(text, opt) //.CreateImageTextLog(text, opt, imgtxtcolor.LogFileAndConsole)
	canvasArr, err := imgtxtcolor.CreateImageTextLog(text, opt, imgtxtcolor.LogFileAndConsole)
	if err != nil {
		return nil, err
	}
	//	dd = imgtxtcolor.GetBase64(canvasArr[0])
	log.Println("Всего картинок:", len(canvasArr))

	return canvasArr, nil
}
