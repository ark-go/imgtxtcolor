package imgtxtcolor

// import (
// 	"image"

// 	"golang.org/x/image/draw"
// )

// func setFrameAroundImg(canvas *ImgCanvas, pathFile string) {
// 	src, left, rigth, top, botttom := getRectInsideFree(pathFile)
// 	w := rigth - left // ширина картинки
// 	pr := float64(src.Rect.Dx()) / float64(w)
// 	fw := float64(canvas.Img.Rect.Dx()) / pr
// 	log.Println("paddingframe", left, rigth, top, botttom)
// 	//src := canvas.Img
// 	//newWidth := canvas.Img.Rect.Dx() + left + rigth
// 	newWidth := int(fw)
// 	newHeight := canvas.Img.Rect.Dy() + top + botttom
// 	imgFrame := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
// 	//imgFrame := image.NewRGBA(image.Rect(0, 0, canvas.Img.Rect.Dx(), canvas.Img.Rect.Dy()))
// 	// Resize:
// 	draw.NearestNeighbor.Scale(imgFrame, imgFrame.Rect, src, src.Bounds(), draw.Over, nil)

// 	draw.Draw(canvas.Img, canvas.Img.Bounds(), imgFrame, image.Point{}, draw.Over)

// }
