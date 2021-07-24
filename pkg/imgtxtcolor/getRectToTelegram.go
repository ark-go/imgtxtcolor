package imgtxtcolor

func getRectToTelegram(width, height float64) (float64, float64) {
	// когдато надо было держать соотношение сторон.. пока отложено но должно работать
	var resHeight float64
	if width >= height {
		cfHeight := width / 2.5
		if cfHeight > height {
			//x := cfHeight - height // сколько не хватает
			resHeight = cfHeight //height + x
		} else {
			resHeight = height
		}
	} else {
		resHeight = height
	}
	return width, resHeight
}
