package imgtxtcolor

func (p *stParam) textToHeight() {
	if len(p.allImages) < 1 {
		return
	}
	switch p.opt.AlignHeight {
	case "center":
		textToCenterHeight(p)
	case "bottom":
		textToBottomHeight(p)
	}
}
