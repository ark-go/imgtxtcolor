package imgtxtcolor

import ()

func drawLine2(param *stParam, word string) bool {
	// новый Image если еще нет или требуют новый
	if param.canvas == nil || param.isNewCanvas {
		if len(word) > 0 || param.canvas == nil {
			param.addNextCanvas()
		}
	}
	// вычислим на сколько продвинется "курсор" по горизонтали
	param.textWidthSumm += param.drw.MeasureString(word)
	// выясним меняется ли высота в строке
	if param.textHeightTmp < param.getHeight() {
		param.textHeightTmp = param.getHeight()
		// TODO
	}

	if ok := checkHeight(param, word); !ok {
		if len(word) > 0 {
			param.addNextCanvas()
			if ok := checkHeight(param, word); !ok {
				return false // нет места для строк
			}
			//			log.Println("height text", param.textHeightSumm.Ceil(), text)
		} else {
			// у нас перебор, но не влезает пустая строка, нечего печатать
			return true
		}
	}

	//log.Println("draw:", param.drw.Dot.Y)
	param.drw.DrawString(word)
	return true
}