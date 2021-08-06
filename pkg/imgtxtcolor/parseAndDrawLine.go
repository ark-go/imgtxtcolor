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
	wordCount := 0
	for word := range c {
		if word == "time:now" {
			t := time.Now()
			word = t.Format("15:04:05")
		}
		if word == "date:now" {
			t := time.Now()
			word = t.Format("02.01.2006")
		}
		// проверка на команду
		if str, cmd := commandGet(word); cmd != "" { // команду пропускаем
			tBreak := false
			if isCmd, tBreak = p.commandCheck(str, cmd); isCmd {
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
			p.textAlignVertical() // перерисовка текста и заявка на новый Image
			isBreak = false
		}
		if p.isNewCanvas {
			p.addNextCanvas()
		}
		wordCount++
		sbTmp.WriteString(word) // temp - только для измерения длинны
		// определим куда переместится курсор
		textWidh := p.drw.MeasureString(sbTmp.String())

		// если вылезает новая временная строка, то пишем предыдущую
		if textWidh.Ceil() > (p.canvas.img.Rect.Dx() - p.canvas.padding.lenW()) {
			if wordCount == 1 {
				// печатаем слова которые не влезают в строку целиком, т.е. с начала строки
				wordCount = 0                 // сброс счетчика слов
				run := []rune(sbTmp.String()) // разбиваем на руны

				for len(run) > 0 { // будем печатать все слово и нарезать его
					i := len(run)
					for ; i > 0; i-- { // поиск наименьшего куска
						rstr := run[:i]                            // отрезаем конец
						textW := p.drw.MeasureString(string(rstr)) // замеряем остаток
						// if (p.canvas.img.Rect.Dx() - p.canvas.padding.lenW() - p.drw.MeasureString("G").Ceil()) < 0 {
						// 	//log.Println(p.canvas.img.Rect.Dx(), p.canvas.padding.lenW(), p.drw.MeasureString("G").Ceil())
						// 	return errors.New("нет места по горизонтали")
						// }
						if textW.Ceil() > (p.canvas.img.Rect.Dx() - p.canvas.padding.lenW()) {
							continue // еще не влезает
						}
						if Ok := p.drawLine(string(rstr)); !Ok { // влезает, печатаем кусок
							return errors.New("нет места по вертикали на новой канве")
						}
						break // напечатали часть,  проверим что осталось
					}
					if i == 0 {
						// поскольку мы дошли до начала, напечатаем и выйдем из for
						// печатать будем даже если и не влезаем!
						p.drawLine(string(run[i:]))
						break
					}
					run = run[i:] // берем остатки, нам их все надо напечатать
				}
				word = "" // закончили печать строки очистим буфер
			} else { // не влезает строка с новым словом, печатаем то что было до слова
				sbStr := strings.TrimRight(sb.String(), " ")

				if Ok := p.drawLine(sbStr); !Ok { // sbStr
					return errors.New("нет места по вертикали на новой канве")
				}
			}
			wordCount = 0
			sb.Reset()
			sbTmp.Reset()
			word = strings.TrimLeft(word, " ") // если переносится пробел впереди - удалим его
			sbTmp.WriteString(word)            // слово которое не влезло вставляем вначало для TMP
			sb.WriteString(word)               // ну и начинаем новый набор с пропущеного слова
			continue
		}

		sb.WriteString(word)
	}
	// если Break был на одельной строке и после него, в строке, не было текста
	if isBreak {
		p.textAlignVertical() // перерисовка текста и заявка на новый Image
		isBreak = false
	}

	sbStr := strings.TrimRight(sb.String(), " ")

	if Ok := p.drawLine(sbStr); !Ok {
		return errors.New("нет места по вертикали на новой канве2")
	}

	return nil
}
