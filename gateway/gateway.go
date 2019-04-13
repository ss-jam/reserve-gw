package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"whiteswan.com/tnsp"
	"whiteswan.com/txsp"
)

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

var mux map[string]func(http.ResponseWriter, *http.Request)
var xhr map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:    ":8000",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = hello
	mux["/tnsp"] = tnsp.Reply
	mux["/txsp"] = txsp.Reply

	xhr = make(map[string]func(http.ResponseWriter, *http.Request))
	xhr["/"] = hello
	xhr["tnsp"] = tnsp.Xhr
	xhr["txsp"] = txsp.Xhr

	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getHandler(w, r)
	case "POST":
		log.Printf("POST Request: %q", r)
	case "PUT":
		log.Printf("PUT Request: %q", r)
	case "DELETE":
		log.Printf("DELETE Request: %q", r)
	default:
		log.Printf("Unhandled Request: %q", r)
	}
	return
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	} else {
		log.Printf("Trying again: %s", r.URL.String())
		p := strings.Split(r.URL.String(), "/")
		log.Printf("Looking for registered root in %q", p)
		if h, ok := mux["/"+p[1]]; ok {
			h(w, r)
			return
		} else if x, ok := r.Header["X-Requested-With"]; ok {
			// Valid Ajax calls may use relative addressing and have a referrer of this server,
			// so validate and fix through matching the referrer
			log.Printf("XHR so trying that...: %q", x)
			// get the entry that matches in mux
			if x[0] == "XMLHttpRequest" {
				ref := strings.Split(r.Referer(), "/")
				log.Printf("Checking referrer: %q", ref)
				for _, item := range ref {
					if h, ok := xhr[item]; ok {
						// This should change the URL to the real relative location and execute the page
						h(w, r)
						return
					}
				}
			}
		} else {
			io.WriteString(w, "My server: "+r.URL.String()+", "+r.Referer())
		}
	}
}
