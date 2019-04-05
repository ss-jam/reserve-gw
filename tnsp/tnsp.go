package tnsp

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"whiteswan.com/rehome"
)

// Return the output to the event trigger
func Reply(w http.ResponseWriter, r *http.Request) {
	//pagefmt := "<html>\n<head>%s\n</head>\n<body>\n<h1>%s</h1>\n<div>\n<div>%s</div>\n</div>\n</body>\n</html>"
	url := "https://reserve.tnstateparks.com"
	resp, err := getSite(url)
	defer resp.Body.Close()
	if err == nil {
		spPage, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			doc := rehome.FixAttributes(spPage, url)
			io.WriteString(w, doc)
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

// This should be moved to its own package for getting reservation and other sites
func getSite(s string) (*http.Response, error) {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	client := &http.Client{
		CheckRedirect: redirectPolicy,
		Transport:     tr,
	}
	resp, err := client.Get(s)
	if err != nil {
		log.Printf("Cannot create new request: %s", err)
		return nil, err
	}

	log.Printf("Cookies: %v\n", client.Jar)
	return resp, nil
}
