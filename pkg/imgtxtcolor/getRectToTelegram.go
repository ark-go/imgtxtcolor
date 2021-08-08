package imgtxtcolor

func getRectToTelegram(width, height float64) (float64, float64) {
	// когдато надо было держать соотношение сторон.. пока отложено но должно работать
	var resHeight float64
	var resWidth float64
	if width >= height {
		cfHeight := width / 2.5
		if cfHeight > height {
			//x := cfHeight - height // сколько не хватает
			resHeight = cfHeight //height + x
		} else {
			resHeight = height
		}
		return width, resHeight
	} else {
		cfWidth := height / 2.5
		if cfWidth < height {
			resWidth = cfWidth
		} else {
			resWidth = width
		}
		return resWidth, height
	}

}
