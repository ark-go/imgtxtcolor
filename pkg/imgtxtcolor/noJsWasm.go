//go:build !js

package imgtxtcolor

import (
	"os"
	"path/filepath"
)

func init() {
	var err error
	rootDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln("не определить рабочий каталог")

	}
}
