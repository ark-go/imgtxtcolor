package imgtxtcolor

import (
	logger "log"
	"os"
)

var log *logger.Logger

func init() {
	log = logger.New(os.Stdout, ": ", logger.LstdFlags)
}
