package internal

import "net/http"

func start(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/test.html")
}
