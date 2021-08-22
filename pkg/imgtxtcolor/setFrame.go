package imgtxtcolor

import (
	"image"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func setFrame(canvas *ImgCanvas, pathFrame string) {

	input, err := os.Open(pathFrame)
	if err != nil {
		log.Println("frame error:", err.Error())
		return
	}
	defer input.Close()

	// output, _ := os.Create("your_image_resized.png")
	// defer output.Close()

	// Decode the image (from PNG to image.Image):
	src, _ := png.Decode(input)

	// Set the expected size that you want:
	//imgFrame := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))
	imgFrame := image.NewRGBA(image.Rect(0, 0, canvas.Img.Rect.Dx(), canvas.Img.Rect.Dy()))

	// Resize:
	draw.NearestNeighbor.Scale(imgFrame, imgFrame.Rect, src, src.Bounds(), draw.Over, nil)

	draw.Draw(canvas.Img, canvas.Img.Bounds(), imgFrame, image.Point{}, draw.Over)
	// Encode to `output`:
	//png.Encode(output, dst)
}

// func loadDecodePNG(p *stParam, pathFrame string) (image.Image, error) {
// 	if fileExists(pathFrame) {
// 		input, err := os.Open(pathFrame)
// 		if err != nil {
// 			log.Println("frame error:", err.Error())
// 			return nil, err
// 		}
// 		defer input.Close()
// 		// Decode the image (from PNG to image.Image):
// 		src, _ := png.Decode(input)
// 		return src, nil
// 	}else{
// 		input, err :=  pathFrame)
// 		if err != nil {
// 			log.Println("frame error:", err.Error())
// 			return nil, err
// 		}
// 		defer input.Close()
// 		// Decode the image (from PNG to image.Image):
// 		src, _ := png.Decode(input)
// 		return src, nil
// 	}
// }

/*
draw.NearestNeighbor
NearestNeighbor - интерполятор ближайшего соседа. Это очень быстро, но обычно дает очень некачественные результаты. При увеличении масштаб результат будет выглядеть «блочным».

draw.ApproxBiLinear
ApproxBiLinear - это смесь интерполяторов ближайшего соседа и билинейных интерполяторов. Это быстро, но обычно дает результаты среднего качества.

draw.BiLinear
BiLinear - это ядро ​​палатки. Это медленно, но обычно дает качественные результаты.

draw.CatmullRom
CatmullRom - это ядро ​​Catmull-Rom. Это очень медленно, но обычно дает очень качественные результаты.
*/
