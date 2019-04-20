package multiplex

import (
	"io"
	"log"
	"net/http"
	"regexp"

	"whiteswan.com/tnsp"
	"whiteswan.com/txsp"
)

type Multiplex struct {
	RefURL  string
	EvalSub bool
	Handler func(http.ResponseWriter, *http.Request)
}

var mux = make(map[string]Multiplex)

// Initialize the system patterns
func Initialize() {
	mux = map[string]Multiplex{
		"/":     {"/", false, hello},
		"/tnsp": {"https://reserve.tnstateparks.com", true, tnsp.Reply},
		"/txsp": {"https://txsp.com", true, txsp.Reply},
	}

	tnsp.Initialize()
}

// Find the Path of a URL string
func ParseUrl(s string) string {
	//r := `[-\w@:%.+~#=]{2,256}\.[a-z]{2,6}\b([-\w@:%+.~#?&//=]*)`
	r := `https?://[\w-@%\.~#+:]{2,256}(/[\w-@%~\.+:]*){1,100}`
	reg := regexp.MustCompile(r)
	l := reg.FindAllStringSubmatch(s, 10)
	log.Printf("Found %q in URL %s", l, s)
	if l != nil {
		return l[0][1]
	} else {
		return ""
	}
}
func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
