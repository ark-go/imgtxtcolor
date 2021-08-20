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
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
