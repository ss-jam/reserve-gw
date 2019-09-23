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

// Add a Multiplex struct of a page reference -
// TODO: add verification of Multiplex values
func AddMux(s string, m Multiplex) bool {
	if mux != nil && mux[s].RefURL == "" && m.RefURL != "" {
		mux[s] = m
		return true
	}
	return false
}

// Get the Multiplex page reference coresponding to the key
// TODO: add return of error or bool along with the mux value
// 	     reflecting if the existence of a Multiplex structure
//       corresponding to the key
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
