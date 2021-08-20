package imgtxtcolor

import (
	"image"

	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func getRectInsideFree(p *stParam, pathFrame string) (img *image.RGBA, yTop, yBottom, xLeft, xRight int) {
	src := getImageFromFile(pathFrame)
	//
	img = image.NewRGBA(image.Rect(0, 0, p.opt.Width, p.opt.Height))
	// Resize:
	draw.NearestNeighbor.Scale(img, img.Rect, src, src.Bounds(), draw.Over, nil)
	//draw.Draw(canvas.Img, canvas.Img.Bounds(), imgFrame, image.Point{}, draw.Over)

	x := img.Rect.Dx() / 2
	y := img.Rect.Dy() / 2
	top := make(chan int)
	bottom := make(chan int)
	left := make(chan int)
	right := make(chan int)
	go getLeftPadding(img, x, y, left)
	go getRightPadding(img, x, y, right)
	go getTopPadding(img, x, y, top)
	go getBottomPadding(img, x, y, bottom)

	for i := 0; i < 4; i++ {
		// ждем два канала
		select {
		case yT := <-top:
			yTop = yT
		case yB := <-bottom:
			yBottom = yB
		case xL := <-left:
			xLeft = xL
		case xR := <-right:
			xRight = xR
		}
	}
	return
}

func getImageFromFile(pathFrame string) *image.RGBA {
	input, err := os.Open(pathFrame)
	if err != nil {
		log.Println("frame error:", err.Error())
		return nil
	}
	defer input.Close()
	// Decode the image (from PNG to image.Image):
	src, _ := png.Decode(input)
	b := src.Bounds()
	imgrgb := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	//draw.Draw(imgrgb, imgrgb.Bounds(), src, b.Min, draw.Src)
	draw.Draw(imgrgb, imgrgb.Bounds(), src, image.Point{}, draw.Src)

	return imgrgb
}

func getLeftPadding(img *image.RGBA, x, y int, c chan int) {
	for i := x; i > img.Bounds().Min.X; i-- { // X left
		//	for j := y; j < y+5; j++ {
		if r, g, b, a := img.At(i, y).RGBA(); (r + g + b + a) != 0 {
			c <- i + 1
			return
		}
		//	}
	}
	c <- 0
}
func getRightPadding(img *image.RGBA, x, y int, c chan int) {
	for i := x; i < img.Bounds().Max.X-1; i++ { // X left
		//	for j := y; j < y+5; j++ {
		if r, g, b, a := img.At(i, y).RGBA(); (r + g + b + a) != 0 {
			c <- i
			return
		}
		//	}
	}
	c <- 0
}
func getTopPadding(img *image.RGBA, x, y int, c chan int) {
	for i := y; i > img.Bounds().Min.Y; i-- { // X left
		if r, g, b, a := img.At(x, i).RGBA(); (r + g + b + a) != 0 {
			c <- i + 1
			return
		}
	}
	c <- 0
}
func getBottomPadding(img *image.RGBA, x, y int, c chan int) {
	for i := y; i < img.Bounds().Max.Y-1; i++ { // X left
		if r, g, b, a := img.At(x, i).RGBA(); (r + g + b + a) != 0 {
			c <- i
			return
		}
	}
	c <- 0
}
