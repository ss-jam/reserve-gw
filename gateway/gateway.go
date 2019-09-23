package main

import (
	"io"
	"log"
	"net/http"

	"whiteswan.com/multiplex"
	"whiteswan.com/remote"
	"whiteswan.com/tnsp"
	"whiteswan.com/txsp"
)

// Initialize the system patterns
func Initialize() {
	remote.InitRemote()

	// Right now, ignore the return value and assume success :)
	multiplex.AddMux("/", multiplex.Multiplex{"/", false, hello})
	multiplex.AddMux("/tnsp", multiplex.Multiplex{"https://reserve.tnstateparks.com", true, tnsp.Reply})
	multiplex.AddMux("/txsp", multiplex.Multiplex{"https://txsp.com", true, txsp.Reply})

	//tnsp.Initialize()
}

// Placeholder function mostly for runtime testing
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	Initialize()

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
