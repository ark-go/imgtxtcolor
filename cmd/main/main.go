package main

import (
	"fmt"
	"os"

	"github.com/ark-go/imgtxtcolor/internal"
)

var versionProg string

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err")
			os.Exit(9)
		}
	}()
	fmt.Println("Версия: ", versionProg)
	internal.StartHttpServer(false)
	//router.Use(middleware.Logger)
	//http.ListenAndServe(":3005", router)
}
