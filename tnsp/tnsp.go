package tnsp

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ss-jam/reserve-gw/remote"

	"github.com/ss-jam/reserve-gw/manipulate"
)

const url = "https://reserve.tnstateparks.com"

var tr *http.Transport
var client *http.Client

// Return the output to the event trigger
func Reply(w http.ResponseWriter, r *http.Request) {
	//pagefmt := "<html>\n<head>%s\n</head>\n<body>\n<h1>%s</h1>\n<div>\n<div>%s</div>\n</div>\n</body>\n</html>"
	log.Printf("TNSP: %s", r.URL.String())
	p := strings.Split(r.URL.Path, "/")
	sel := ""
	if p[1] == "tnsp" && len(p) > 2 {
		sel = r.URL.Path[5:]
	}
	var resp *http.Response
	var err error
	switch r.Method {
	case "GET":
		resp, err = remote.GetRemote(url+sel, r.Header, r.Method)
		//resp, err = getSite(url + sel)
		if err != nil {
			log.Printf("ERROR getting remote: %s", err)
			http.Error(w, err.Error(), resp.StatusCode)
			return
		}
		err = remote.Write(w, resp)
		if err != nil {
			log.Printf("ERROR: could not write response: %s", err)
			return
		}
	case "POST":
		err := r.ParseForm()
		if err != nil {
			log.Printf("Error tyring to parse form: %s", err)
		}
		log.Printf("Query: %q", r.PostForm)
		resp, err = postSite(url+sel, r.Header["Content-Type"][0], strings.NewReader(r.PostForm.Encode()))
		defer resp.Body.Close()
		if err == nil {
			spPage, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				doc, l := manipulate.SimpleBatch(spPage, url, "tnsp")
				io.WriteString(w, doc)
				log.Printf("Changed links: %v", l)
				//doc, err := elements.BatchParse(spPage, url)
				//if err == nil {
				//	page := fmt.Sprintf(pagefmt, doc.Head,
				//		"Tennessee State Parks alternative resource", doc.Body)
				//	io.WriteString(w, page)
				//}
			} else {
				log.Printf("Could not read body: %s", err)
			}
		} else {
			io.WriteString(w, "Welcome to the Tennessee State Parks alternative booking resource.")
		}
	default:
		log.Printf("Unimplemented method: %s", r.Method)
		http.Error(w, "Something unexpected this way comes", 404)
		return
	}
}

// Return the contents of a relative XHR request
func Xhr(w http.ResponseWriter, r *http.Request) {
	log.Printf("XHR URL: %s", r.URL.String())
	resp, err := getSite(url + r.URL.String())
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, url+r.URL.String()+": Not Found", 404)
		return
	}
	spPage, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Could not read XHR body: %s", err)
		http.Error(w, url+r.URL.String()+": Not Found", 404)
		return
	}
	doc, _ := manipulate.SimpleBatch(spPage, url, "tnsp")
	io.WriteString(w, doc)
}

// This is here just for investigation/learning
func redirectPolicy(req *http.Request, via []*http.Request) error {
	if len(via) > 0 {
		log.Printf("Past requests (%d):\n", len(via))
		for i, v := range via {
			log.Printf("%d: %s\n", i, v.URL)
		}
	} else {
		log.Print("No redirections\n")
	}
	return nil
}

func Initialize() {
	tr = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client = &http.Client{
		CheckRedirect: redirectPolicy,
		Transport:     tr,
	}
}

// This should be moved to its own package for getting reservation and other sites
func getSite(s string) (*http.Response, error) {
	resp, err := client.Get(s)
	if err != nil {
		log.Printf("Cannot create new request: %s", err)
		return nil, err
	}

	log.Printf("Cookies: %v\n", client.Jar)
	return resp, nil
}

func postSite(u string, c string, b io.Reader) (*http.Response, error) {
	log.Printf("POST URL: %s, Content-Type: %s", u, c)
	resp, err := client.Post(u, c, b)
	if err != nil {
		log.Printf("Cannot create new request: %s", err)
		return nil, err
	}

	log.Printf("POST SITE: %q", resp)
	return resp, nil
}
