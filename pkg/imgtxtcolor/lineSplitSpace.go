package imgtxtcolor

import "strings"

// разделяет строку по словам, по первому пробелу
// 	line - строка
// 	def - это пробел или символ заменяющий пробел
// 	с - канал передачи результата
func lineSplitSpace(line string, def string, c chan string) {
	f := strings.Split(line, " ") // пробелы исчезают
	for i := 0; i < len(f); i++ {
		if f[i] == "" { // вероятно был лишний пробел, восстановим
			c <- def
		} else {
			c <- f[i]
		}
		if f[i] != "" && i != len(f)-1 { // добавляем пробел если не последнее слово
			//это пробел на котором режется слово, поэтому тут добавим
			c <- def
		}
	}
	close(c)
}

/*
c := make(chan string)
	go lineSplitSpace("1 234  567890", "+", c)
	for v := range c {
		fmt.Printf("%s\n", v)
	}
*/
