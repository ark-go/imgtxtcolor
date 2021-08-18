package internal

import (
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

type ViewData struct {
	SizeW        int
	SizeH        int
	Initials     string
	FontSize     int
	FimagNames   []string
	FimageBase64 []string
}

func getParamInt(r *http.Request, nameKey string, def int) int {
	val, err := strconv.Atoi(r.FormValue(nameKey))
	if err != nil {
		val = def
	}
	return val
}

func sendImages(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос..")
	initials := r.FormValue("initials")
	sizeH := getParamInt(r, "sizeH", 300)
	sizeW := getParamInt(r, "sizeW", 500)
	fontSizeInt := getParamInt(r, "fontSizeInt", 20)
	fimgNames := []string{}
	var img64 []string
	os.MkdirAll("internal/img/", os.ModePerm)
	avatar, err := createImages(sizeH, sizeW, fontSizeInt, initials)
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
		// отправка base64 вместо png или то и это
		// img64 = getBase64(avatar)
		// log.Println("bas64", "готов")
		// ----- http2  push -----------
		if pusher, ok := w.(http.Pusher); ok {
			// Push is supported.
			options := &http.PushOptions{
				Header: http.Header{
					"Accept-Encoding": r.Header["Accept-Encoding"],
				},
			}
			log.Println("Есть http2/Push")

			// ---------  png files  ------------------
			for i, imgCanvas := range avatar {
				name := "img" + strconv.Itoa(i) + ".png"
				f, err := os.Create("internal/img/" + name)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				png.Encode(f, imgCanvas.Img)
				if err := pusher.Push("/img/"+name, options); err != nil {
					log.Printf("%s Failed to push: %v", name, err)
				}
				fimgNames = append(fimgNames, name)
			}
			//time.Sleep(time.Second * 4)
			// ---------- end png files -----------

		} else {
			// ---------  png files  ------------------
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
	}
	data := ViewData{
		SizeW:        sizeW,
		SizeH:        sizeH,
		Initials:     initials,
		FontSize:     fontSizeInt,
		FimagNames:   fimgNames,
		FimageBase64: img64,
	}
	tmpl, _ := template.ParseFiles("internal/test.html")
	tmpl.Execute(w, data)
	//http.ServeFile(w, r, "internal/test.html")
	log.Println("Слушаем порт 3005  https://127.0.0.1:3005/avatar")
}
