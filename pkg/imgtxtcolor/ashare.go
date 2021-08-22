package imgtxtcolor

import (
	logger "log"
	"os"
)

var log *logger.Logger

func init() {
	log = logger.New(os.Stdout, ": ", logger.LstdFlags)
}

func fileExists(filename string) bool {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error FileExists", filename)

		}
	}()
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false // уверенно нет такого файла или каталога
		}
		return false // хз ошибка какаято но не про файл
	}
	return !info.IsDir()
}
