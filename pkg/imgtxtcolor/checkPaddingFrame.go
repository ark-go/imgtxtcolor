package imgtxtcolor

func checkPaddingFrame(p *stParam) {
	if p.opt.FrameFilePath == "" {
		return
	}
	img, top, bottom, left, right := getRectInsideFree(p, p.opt.FrameFilePath)
	p.opt.Padding.top = top
	p.opt.Padding.bottom = img.Rect.Dy() - bottom
	p.opt.Padding.left = left
	p.opt.Padding.right = img.Rect.Dx() - right
}
