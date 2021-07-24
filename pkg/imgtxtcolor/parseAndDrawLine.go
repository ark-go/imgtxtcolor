package imgtxtcolor

import (
	"errors"
	"strings"
	"time"
)

func (p *stParam) parseAndDrawLine(text string) error {
	//defer duration(track("строка"))
	var sb, sbTmp strings.Builder
	c := make(chan string)
	// запрос на разбивку по словам и пробелам, пробелы тоже есть
	// перебираем слова из канала
	go lineSplitSpace(text, " ", c)

	isCmd := false
	for word := range c {
		if word == "time:now" {
			t := time.Now()
			word = t.Format("15:04:05")
		}
		// проверка на команду
		if str, cmd := commandGet(word); cmd != "" { // команду пропускаем
			if isCmd = commandCheck(p, str, cmd); isCmd {
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
		test := 1
		if test == 1 {
			sbTmp.WriteString(word) // temp - только для измерения длинны
			if p.isNewCanvas {
				p.addNextCanvas()
			}
			// определим куда переместится курсор
			textWidh := p.drw.MeasureString(sbTmp.String())

			// если вылезает новая временная строка, то пишем предыдущую
			if textWidh.Ceil() > (p.canvas.Rect.Dx() - p.padding.lenW()) {
				sbStr := strings.TrimRight(sb.String(), " ")

				if Ok := p.drawLine(sbStr); !Ok {
					return errors.New("нет места по вертикали на новой канве")
				}
				sb.Reset()
				sbTmp.Reset()
				sbTmp.WriteString(word) // слово которое не влезло вставляем вначало для TMP
				sb.WriteString(word)    // ну и начинаем новый набор с пропущеного слова
				continue
			}
		} else if test == 2 {

			if Ok := drawLine2(p, word); !Ok {
				return errors.New("нет места по вертикали на новой канве")
			}
		}
		sb.WriteString(word)
	}

	sbStr := strings.TrimRight(sb.String(), " ")

	if Ok := p.drawLine(sbStr); !Ok {
		return errors.New("нет места по вертикали на новой канве2")
	}

	//	}
	//}
	return nil
}
