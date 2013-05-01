package main

import (
	"github.com/gorilla/mux"
	"github.com/jkallhoff/gofig"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")

	http.Handle("/", router)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	panic(http.ListenAndServe(gofig.Str("webPort"), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/views/home.html")
	return
}
