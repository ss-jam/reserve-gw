package multiplex

import (
	"log"
	"net/http"
	"regexp"
)

// Route pattern matching structure for server
type Multiplex struct {
	RefURL  string
	EvalSub bool
	Handler func(http.ResponseWriter, *http.Request)
}

var mux = make(map[string]Multiplex)

func AddMux(s string, m Multiplex) {
	mux[s] = m
}

func GetMux(s string) Multiplex {
	return mux[s]
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
