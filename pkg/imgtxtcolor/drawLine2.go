package imgtxtcolor

import "golang.org/x/image/math/fixed"

func drawLine2(param *stParam, word string) bool {
	// новый Image если еще нет или требуют новый
	if param.canvas.Img == nil || param.isNewCanvas {
		if len(word) > 0 || param.canvas.Img == nil {
			param.addNextCanvas()
		}
	}

	// выясним меняется ли высота в строке
	metric := param.drw.Face.Metrics()
	if param.textHeightTmp < metric.Height+fixed.I(param.opt.LineSpacing) {
		param.textHeightTmp = metric.Height + fixed.I(param.opt.LineSpacing)
		// TODO
	}

	if ok := param.checkHeight(word); !ok {
		if len(word) > 0 {
			param.addNextCanvas()
			if ok := param.checkHeight(word); !ok {
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
