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
	isBreak := false
	for word := range c {
		if word == "time:now" {
			t := time.Now()
			word = t.Format("15:04:05")
		}
		// проверка на команду
		if str, cmd := commandGet(word); cmd != "" { // команду пропускаем
			tBreak := false
			if isCmd, tBreak = commandCheck(p, str, cmd); isCmd {
				if tBreak {
					isBreak = true
				}
				// только если распознали команду пропускаем и не печатаем слово
				continue
			}
		}
		if isCmd {
			// предыдущим словом была команда, за ней убираем все последующие пробелы
			// word = strings.TrimLeft(word, " ") // удалит пробелы из начала слова
			// if word == "" {
			if strings.TrimLeft(word, " ") == "" {
				// не будем снимать флаг isCmd, а будем удалять всё пустое простраство слева от следущего слова после команды
				continue
			}
			isCmd = false //  снимаем флаг и  разрешаем печатать слово

		}

		if isBreak {
			p.textToHeight() // перерисовка текста и заявка на новый Image
			isBreak = false
		}
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

		sb.WriteString(word)
	}
	// если Break был на одельной строке и после него, в строке, не было текста
	if isBreak {
		p.textToHeight() // перерисовка текста и заявка на новый Image
		isBreak = false
	}

	sbStr := strings.TrimRight(sb.String(), " ")

	if Ok := p.drawLine(sbStr); !Ok {
		return errors.New("нет места по вертикали на новой канве2")
	}

	return nil
}
