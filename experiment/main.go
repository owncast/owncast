package main

import (
	"net/http"
)

func main() {
	s := Handlers{}
	h := Handler(s)

	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/welcome", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// })
	// r.Mount("/", Handler(&handlers))
	http.ListenAndServe(":3000", h)
}
