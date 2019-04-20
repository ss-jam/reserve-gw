package main

import (
	"log"
	"net/http"

	"whiteswan.com/multiplex"
)

func main() {
	multiplex.Initialize()

	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		multiplex.GetHandler(w, r)
	case "POST":
		multiplex.PostHandler(w, r)
	case "PUT":
		log.Printf("PUT Request: %q", r)
	case "DELETE":
		log.Printf("DELETE Request: %q", r)
	default:
		log.Printf("Unhandled Request: %q", r)
	}
	return
}
