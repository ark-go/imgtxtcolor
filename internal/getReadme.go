package internal

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"io/ioutil"
	"log"
	"net/http"
)

func getReadme() []byte {
	//resp, err := http.Get("https://github.com/ark-go/imgtxtcolor/blob/main/README.md")
	resp, err := http.Get("https://raw.githubusercontent.com/ark-go/imgtxtcolor/main/README.md")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	//md := []byte("markdown text")
	html := markdown.ToHTML(body, parser, nil)
	return html
}
