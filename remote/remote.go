package remote

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Create a transport and a client for retrieval of remote content
var tr *http.Transport
var client *http.Client

func InitRemote() {
	tr = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client = &http.Client{
		Transport: tr,
	}
}

// Used for local referrer relative links of remote content
func GetRemote(s string, h http.Header) (*http.Response, error) {
	req, err := http.NewRequest("GET", s, nil)
	if err != nil {
		log.Printf("Error constructing request: %s", err)
		return nil, err
	}
	//log.Printf("getting %s with type: %s", s, h.Get("Content-Type"))
	req.Header = h //.Set("Content-Type", h.Get("Content-Type"))
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error getting remote: %s", err)
		return nil, err
	}
	return resp, nil
}

// Write out the contents retrieved from GetRemote or other http.Response
func Write(w http.ResponseWriter, r *http.Response) error {
	defer r.Body.Close()
	bod, err := ioutil.ReadAll(r.Body)
	log.Printf("Resp Header (%s): %q", r.Request.Host, r.Header)
	for k, vs := range r.Header {
		for _, v := range vs {
			w.Header().Set(k, v)
		}
	}
	if err != nil {
		http.Error(w, err.Error(), r.StatusCode)
		return err
	}
	w.WriteHeader(r.StatusCode)
	w.Write(bod)
	return nil
}
