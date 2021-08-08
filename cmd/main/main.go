package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"path"

	"image/png"
	"net/http"
	"os"
	"strconv"

	"github.com/ark-go/imgtxtcolor/pkg/imgtxtcolor"
	"github.com/go-chi/chi"
	_ "github.com/go-chi/chi/middleware"
)

type ViewData struct {
	SizeW      int
	SizeH      int
	Initials   string
	FontSize   int
	FimagNames []string
}

func main() {
	router := chi.NewRouter()
	//router.Use(middleware.Logger)
	router.Get("/start", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/test.html")
	})
	router.Get("/avatar", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Запрос..")
		initials := r.FormValue("initials")
		sizeH, err := strconv.Atoi(r.FormValue("sizeH"))
		if err != nil {
			sizeH = 200
		}
		sizeW, err := strconv.Atoi(r.FormValue("sizeW"))
		if err != nil {
			sizeW = 200
		}
		fontSizeInt, err := strconv.Atoi(r.FormValue("fontSizeInt"))
		if err != nil {
			sizeW = 200
		}
		fimgNames := []string{}
		os.MkdirAll("internal/img/", os.ModePerm)
		avatar, err := createAvatar(sizeH, sizeW, fontSizeInt, initials)
		if err != nil {
			input, err := ioutil.ReadFile("internal/errormsg/error.png")
			if err != nil {
				log.Println(err)
				return
			}

			err = ioutil.WriteFile("internal/img/error.png", input, 0644)
			if err != nil {
				log.Println("Error creating", err.Error())
				return
			}
			fimgNames = append(fimgNames, "error.png")
		} else {

			for i, imgCanvas := range avatar {
				name := "img" + strconv.Itoa(i) + ".png"
				f, err := os.Create("internal/img/" + name)
				if err != nil {
					panic(err)
				}
				defer f.Close()

				png.Encode(f, imgCanvas.Img)
				fimgNames = append(fimgNames, name)
			}
		}
		data := ViewData{
			SizeW:      sizeW,
			SizeH:      sizeH,
			Initials:   initials,
			FontSize:   fontSizeInt,
			FimagNames: fimgNames,
		}
		tmpl, _ := template.ParseFiles("internal/test.html")
		tmpl.Execute(w, data)
		//http.ServeFile(w, r, "internal/test.html")
	})
	// router.Get("/img.png", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "internal/img/img.png")
	// })
	router.Get("/img/{name}", func(w http.ResponseWriter, r *http.Request) {
		val := chi.URLParam(r, "name")
		http.ServeFile(w, r, "internal/img/"+val)
	})

	log.Println("Старт listen")
	http.ListenAndServe(":3005", router)
}
func createAvatar(sizeH, sizeW, fontSizeInt int, text string) ([]*imgtxtcolor.ImgCanvas, error) {
	dir, _ := ioutil.ReadDir("internal/img")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"internal/img", d.Name()}...))
	}
	opt := imgtxtcolor.StartOption()
	opt.Width = sizeW
	opt.Height = sizeH
	opt.FontSize = fontSizeInt
	opt.GifFileName = "internal/img/test.gif"
	opt.GifDelay = 100 * 1

	canvasArr, err := imgtxtcolor.CreateImageTextLog(text, opt, imgtxtcolor.LogFileAndConsole)

	if err != nil {
		return nil, err
	}
	//	dd = imgtxtcolor.GetBase64(canvasArr[0])
	log.Println("Всего картинок:", len(canvasArr))

	return canvasArr, nil
}
