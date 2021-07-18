package imgtxtcolor

import "strings"

// Проверка и выделение команды, если есть
// 	возврат:
// 	str - либо параметр команды, либо возврат входной строки
// 	cmd - команда вернется только если она есть
func commandGet(word string) (str string, cmd string) {
	f := strings.SplitN(word, ":", 2)
	if len(f) == 1 {
		return word, ""
	}
	if len(f) == 2 {
		if f[0] == "" || f[1] == "" {
			return word, ""
		}
		// TODO Проверку на команду сюда
		return f[1], f[0]
	}
	return word, "" // этого не может быть
}
