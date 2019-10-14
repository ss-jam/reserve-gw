package multiplex

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"

	"reserve-gw/remote"
)

func getBase(s string) string {
	parts := strings.Split(s, "/")
	base := ""
	for _, part := range parts {
		// Relative reference
		if _, ok := mux["/"+part]; len(part) > 0 && ok {
			base = "/" + part
			break
		}
	}
	return base
}

func makeURL(r *http.Request) (string, error) {
	base := getBase(r.Referer())
	if base == "" {
		return base, errors.New("Could not find base URL reference")
	}
	log.Printf("Attempting to get relative link data for %s using %s", r.URL.String(), base)
	//t := ParseUrl(r.Referer())
	//log.Printf("TEST: %s", t)
	newr := mux[base].RefURL + r.URL.Path
	log.Printf("URL: %s", newr)
	return newr, nil
}

// GetHandler is The GET method handler
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h.Handler(w, r)
		return
	}

	// Looking for any sub references as in /hello/*
	// To be more secure, this should be assigned to
	// a key in the structure instead of always processed.
	log.Printf("Trying again: %s (PATH: %s [QUERY: %q, HEADER: %q])", r.URL.String(), r.URL.Path, r.URL.Query(), r.Header)
	p := strings.Split(r.URL.Path, "/")
	log.Printf("Looking for registered root in %q", p)
	if h, ok := mux["/"+p[1]]; ok && h.EvalSub {
		h.Handler(w, r)
		return
	} else if x, ok := r.Header["X-Requested-With"]; ok {
		// Valid Ajax calls may use relative addressing and have a referrer of this server,
		// so validate and fix through matching the referrer
		log.Printf("XHR so trying that...: %q", x)
		// get the entry that matches in mux
		if x[0] == "XMLHttpRequest" {
			u, err := url.Parse(r.Referer())
			if err != nil {
				http.Error(w, "Malformed referrer: '"+r.Referer()+"'", 404)
				return
			}
			ref := strings.Split(u.Path, "/")
			log.Printf("Checking referrer: %q", ref)
			for _, item := range ref {
				if h, ok := mux["/"+item]; ok {
					// This should change the URL to the real relative location and execute the page
					h.Handler(w, r)
					return
				}
			}
		}
	} else {
		newr, err := makeURL(r)
		if err != nil {
			http.Error(w, "My server (GET): "+r.URL.String()+" ("+r.URL.Path+"), "+r.Referer(), 404)
			return
		}
		resp, err := remote.GetRemote(newr, r.Header, "GET")
		if err != nil {
			log.Printf("Error getting link: %s", err)
			return
		}
		err = remote.Write(w, resp)
		if err != nil {
			log.Printf("ERROR writing response: %s", err)
			return
		}
	}
}

// PostHandler is The POST method handler
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h.Handler(w, r)
		return
	}

	log.Printf("Trying again: %s", r.URL.String())
	p := strings.Split(r.URL.Path, "/")
	log.Printf("Looking for registered root in %q", p)
	if h, ok := mux["/"+p[1]]; ok {
		h.Handler(w, r)
		return
	}

	if x, ok := r.Header["X-Requested-With"]; ok {
		// Valid Ajax calls may use relative addressing and have a referrer of this server,
		// so validate and fix through matching the referrer
		log.Printf("XHR so trying that...: %q", x)
		// get the entry that matches in mux
		if x[0] == "XMLHttpRequest" {
			ref := strings.Split(r.Referer(), "/")
			log.Printf("Checking referrer: %q", ref)
			for _, item := range ref {
				if h, ok := mux["/"+item]; ok {
					// This should change the URL to the real relative location and execute the page
					h.Handler(w, r)
					return
				}
			}
		}
	} else {
		newr, err := makeURL(r)
		if err != nil {
			http.Error(w, "My server (POST): "+r.URL.String()+" ("+r.URL.Path+"), "+r.Referer(), 404)
			return
		}
		resp, err := remote.GetRemote(newr, r.Header, "POST")
		if err != nil {
			log.Printf("Error getting link: %s", err)
			return
		}
		err = remote.Write(w, resp)
		if err != nil {
			log.Printf("ERROR writing response: %s", err)
			return
		}
	}
}
