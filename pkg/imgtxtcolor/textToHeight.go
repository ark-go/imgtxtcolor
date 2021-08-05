package imgtxtcolor

func (p *stParam) textToHeight() {
	if len(p.allImages) < 1 {
		return
	}
	switch p.canvasOpt.alignVertical {
	case AlignVerticalTop:
		//textToCenterHeight(p)
		//	p.isNewCanvas = true // новый текст в новом Image
		textToTopHeight(p)
	case AlignVerticalCenter:
		textToCenterHeight(p)
	case AlignVerticalBottom:
		textToBottomHeight(p)
	}
}
