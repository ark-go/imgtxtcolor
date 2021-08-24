package internal

import (
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

//go:embed test.html
//go:embed errormsg/error.png
//go:embed testembed.txt
//go:embed frameEmbed
var embedFS embed.FS
var rootDir string
var frameDir string
var imgDir string

func init() {
	var err error
	rootDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	imgDir = filepath.Join(rootDir, "img")
	if err != nil {
		rootDir = ""
		log.Println("не определить рабочий каталог")
		return
	}
	getExapleFrame()

}

func getExapleFrame() {
	frameDir = filepath.Join(rootDir, "frame")
	os.MkdirAll(frameDir, os.ModePerm)

	frames, _ := embedFS.ReadDir("frameEmbed")
	for _, file := range frames {

		//	data, _ := embedFS.ReadFile(filepath.Join("frameEmbed", file.Name()))
		data, _ := embedFS.ReadFile("frameEmbed/" + file.Name())

		//	fl, _ := embedFS.Open(filepath.Join(frameDir, file.Name()))
		//	data, _ := fl.Read()
		//fmt.Println(file.Name())
		ioutil.WriteFile(filepath.Join(frameDir, file.Name()), data, 0644)
	}
	fmt.Println("Созданы примеры рамок")
}

// func start(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "internal/test.html")
// }

// func copy(src string, dst string) {
// 	data, _ := ioutil.ReadFile(src)
// 	ioutil.WriteFile(dst, data, 0644)
// }
