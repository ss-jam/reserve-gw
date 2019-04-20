package multiplex

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"whiteswan.com/tnsp"
)

// The GET method handler
func GetHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.Referer(), "/")
	base := ""
	for _, part := range parts {
		// Relative reference
		if _, ok := mux["/"+part]; len(part) > 0 && ok {
			base = "/" + part
			break
		}
	}
	if h, ok := mux[r.URL.String()]; ok {
		h.Handler(w, r)
		return
	} else {
		// Looking for any sub references as in /hello/*
		// To be more secure, this should be assigned to
		// a key in the structure instead of always processed.
		log.Printf("Trying again: %s (PATH: %s [QUERY: %q])", r.URL.String(), r.URL.Path, r.URL.Query())
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
			if base == "" {
				http.Error(w, "My server: "+r.URL.String()+" ("+r.URL.Path+"), "+r.Referer(), 404)
				return
			}
			log.Printf("Attempting to get relative link data for %s using %s", r.URL.String(), base)
			//t := ParseUrl(r.Referer())
			//log.Printf("TEST: %s", t)
			newr := mux[base].RefURL + r.URL.Path
			log.Printf("URL: %s", newr)
			resp, err := tnsp.GetRemote(newr)
			if err != nil {
				log.Printf("Error getting link: %s", err)
				return
			}
			defer resp.Body.Close()
			bod, err := ioutil.ReadAll(resp.Body)
			log.Printf("Resp Header: %q", resp.Header)
			for k, vs := range resp.Header {
				for _, v := range vs {
					w.Header().Set(k, v)
				}
			}
			if err != nil {
				http.Error(w, err.Error(), resp.StatusCode)
				return
			}
			w.WriteHeader(resp.StatusCode)
			w.Write(bod)
		}
	}
}

// The POST method handler
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if h, ok := mux[r.URL.String()]; ok {
		h.Handler(w, r)
		return
	} else {
		log.Printf("Trying again: %s", r.URL.String())
		p := strings.Split(r.URL.Path, "/")
		log.Printf("Looking for registered root in %q", p)
		if h, ok := mux["/"+p[1]]; ok {
			h.Handler(w, r)
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
					if h, ok := mux["/"+item]; ok {
						// This should change the URL to the real relative location and execute the page
						h.Handler(w, r)
						return
					}
				}
			}
		} else {
			io.WriteString(w, "My server: "+r.URL.String()+", "+r.Referer())
		}
	}
}
