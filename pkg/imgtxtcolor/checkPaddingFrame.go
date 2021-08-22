package imgtxtcolor

func checkPaddingFrame(p *stParam) {
	if p.opt.FrameFilePath == "" {
		return
	}
	img, top, bottom, left, right := getRectInsideFree(p, p.opt.FrameFilePath)
	if img == nil {
		p.opt.FrameFilePath = "" // ошибка файла, сбросим путь т.к. он является флагом в дальнейшем
		return
	}
	///   if (top+bottom) > p.canvas.Img.Rect.Dy()  || (left + right) > p.canvas.im

	p.opt.Padding.top = top
	p.opt.Padding.bottom = bottom
	p.opt.Padding.left = left
	p.opt.Padding.right = right
}
