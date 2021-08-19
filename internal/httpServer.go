package internal

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"golang.org/x/net/http2"
)

var router = chi.NewRouter()

func StartHttpServer(https bool) {
	router.Use(redirectHttps)
	router.Get("/start", start)
	router.Get("/avatar", sendImages)
	router.Post("/avatar", sendImages)
	router.Get("/img/{name}", func(w http.ResponseWriter, r *http.Request) {
		val := chi.URLParam(r, "name")
		http.ServeFile(w, r, "internal/img/"+val)
	})
	router.Get("/favicon.ico", sendFavicon)
	if https {
		http2serv()
	} else {
		httpserv()
	}

}

func http2serv() {
	httpServer := http.Server{
		Addr:         ":3005",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
		TLSConfig:    tlsConfig(),
	}
	//	go redirectToHTTPS(httpServer.Addr)
	var http2Server = http2.Server{}
	if err := http2.ConfigureServer(&httpServer, &http2Server); err != nil {
		log.Println("Ошибка http2")
	}
	log.Println("Старт... слушаем порт 3005  https://127.0.0.1:3005/avatar")
	err := httpServer.ListenAndServeTLS("", "")
	if err != nil {
		log.Println("start server https Error:", err.Error())
	}
}
func httpserv() {
	httpServer := http.Server{
		Addr:         ":3005",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	log.Println("Старт... слушаем порт 3005  http://127.0.0.1:3005/avatar")
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Println("start server http Error:", err.Error())
	}
}

func redirectHttps(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	log.Println("middleware...", r.URL.Path)
		// for key, val := range r.Header {
		// 	log.Println(key, val)
		// }
		// ctx := context.WithValue(r.Context(), "article", article)
		// next.ServeHTTP(w, r.WithContext(ctx))
		next.ServeHTTP(w, r)
	})
}

// func redirectToHTTPS(httpAddr string) {
// 	httpSrv := http.Server{
// 		Addr: httpAddr,
// 		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			//	host, _, _ := net.SplitHostPort(r.Host)
// 			u := r.URL
// 			//	u.Host = net.JoinHostPort(host, tlsPort)
// 			u.Scheme = "https"
// 			log.Println(u.String())
// 			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
// 		}),
// 	}
// 	log.Println(httpSrv.ListenAndServe())
// }
