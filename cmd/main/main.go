package main

import (
	"fmt"
	"html/template"
	"image"
	"io/ioutil"
	"path"

	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/ark-go/canvas/pkg/imgtxtcolor"
	"github.com/go-chi/chi"
	_ "github.com/go-chi/chi/middleware"
)

type ViewData struct {
	SizeW       int
	SizeH       int
	Initials    string
	FontSizeInt int
	FimagNames  []string
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
				fmt.Println(err)
				return
			}

			err = ioutil.WriteFile("internal/img/error.png", input, 0644)
			if err != nil {
				fmt.Println("Error creating", err.Error())
				return
			}
			fimgNames = append(fimgNames, "error.png")
		} else {

			for i, img := range avatar {
				_ = img
				name := "img" + fmt.Sprintf("%d", i) + ".png"
				f, err := os.Create("internal/img/" + name)
				if err != nil {
					panic(err)
				}
				defer f.Close()

				png.Encode(f, img)
				fimgNames = append(fimgNames, name)
			}
		}
		data := ViewData{
			SizeW:       sizeW,
			SizeH:       sizeH,
			Initials:    initials,
			FontSizeInt: fontSizeInt,
			FimagNames:  fimgNames,
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
func createAvatar(sizeH, sizeW, fontSizeInt int, text string) ([]*image.RGBA, error) {
	dir, _ := ioutil.ReadDir("internal/img")
	for _, d := range dir {
		os.RemoveAll(path.Join([]string{"internal/img", d.Name()}...))
	}
	opt := imgtxtcolor.StartOption()
	opt.Width = sizeW
	opt.Height = sizeH
	opt.FontSizeInt = fontSizeInt

	canvasArr, err := imgtxtcolor.CreateImageText(text, opt)
	if err != nil {
		return nil, err
	}
	//	dd = imgtxtcolor.GetBase64(canvasArr[0])
	log.Println("Всего картинок:", len(canvasArr))
	return canvasArr, nil
}
