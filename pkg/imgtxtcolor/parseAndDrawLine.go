package imgtxtcolor

import (
	"errors"
	"strings"
)

func parseAndDrawLine(param *stParam, text string) error {
	//defer duration(track("строка"))
	var sb, sbTmp strings.Builder
	c := make(chan string)
	// defer func() {
	// 	select {
	// 	case <-c:
	// 	default:
	// 		//"Channel is not closed"
	// 		close(c)
	// 	}
	// }()

	// запрос на разбивку по словам и пробелам, пробелы тоже есть
	// перебираем слова из канала
	go lineSplitSpace(text, " ", c)

	//var tmpHeight fixed.Int26_6
	isCmd := false
	isPageBreak := false
	for word := range c {

		// проверка на команду
		if str, cmd := commandGet(word); cmd != "" { // команду пропускаем
			if isPageBreak = commandCheckBreak(param, str, cmd); isPageBreak {
				// только если распознали команду прерываем
				//	param.addNextCanvas() // новая канва
				//param.textHeightSumm = fixed.I(0)
				param.isNewCanvas = true
				isCmd = true
				continue
			}
			if isCmd = commandCheck(param, str, cmd); isCmd {
				// только если распознали команду пропускаем и не печатаем слово
				continue
			}
		}
		if isCmd {
			// предыдущим словом была команда, или Break  за ней убираем все последующие пробелы
			word = strings.TrimLeft(word, " ") // удалит пробелы из начала слова
			if word == "" {
				// не будем снимать флаг isCmd, а будем удалять всё пустое простраство слева от следущего слова после команды
				continue
			}
			isCmd = false //  снимаем флаг и  разрешаем печатать слово

		}

		sbTmp.WriteString(word) // temp - только для измерения длинны
		if param.canvas == nil {
			param.addNextCanvas()
		}
		textWidh := param.drw.MeasureString(sbTmp.String())
		// если вылезает новая временная строка, то пишем предыдущую
		if textWidh.Ceil() > (param.canvas.Rect.Dx() - param.padding.lenW()) {
			sbStr := strings.TrimRight(sb.String(), " ")

			if Ok := drawLine(param, sbStr); !Ok {
				return errors.New("нет места по вертикали на новой канве")
			}
			sb.Reset()
			sbTmp.Reset()
			sbTmp.WriteString(word) // слово которое не влезло вставляем вначало для TMP
			sb.WriteString(word)    // ну и начинаем новый набор с пропущеного слова
			continue
		}

		sb.WriteString(word)
	}

	sbStr := strings.TrimRight(sb.String(), " ")

	//if len(sbStr) > 0 {
	//param.drw.Dot.X = param.xPositionFunc(sbStr)
	//	if Ok := drawLine(param, sbStr); !Ok {
	// не было места
	//		param.addNextCanvas()

	if Ok := drawLine(param, sbStr); !Ok {
		return errors.New("нет места по вертикали на новой канве2")
	}

	//	}
	//}
	return nil
}
