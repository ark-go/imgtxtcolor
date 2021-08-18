package main

import (
	"github.com/ark-go/imgtxtcolor/internal"
)

func main() {
	internal.StartHttpServer(true)
	//router.Use(middleware.Logger)
	//http.ListenAndServe(":3005", router)
}
