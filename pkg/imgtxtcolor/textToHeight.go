package imgtxtcolor

func (p *stParam) textToHeight() {
	if len(p.allImages) < 1 {
		return
	}
	switch p.opt.AlignHeight {
	case "top":
		//textToCenterHeight(p)
		//	p.isNewCanvas = true // новый текст в новом Image
		textToTopHeight(p)
	case "center":
		textToCenterHeight(p)
	case "bottom":
		textToBottomHeight(p)
	}
}
