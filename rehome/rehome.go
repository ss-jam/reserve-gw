package rehome

import (
	"bytes"
	"fmt"
	"log"
	"strings"
)

// fix any src or href attributes that are relative to the reference URL
func FixAttributes(b []byte, url string) string {
	//	h := strings.Index(s, "<head")
	//	he := strings.Index(s, "</head>")
	//	b := strings.Index(s, "<body")
	//	be := strings.Index(s, "</body>")
	s := string(b[:])
	//log.Printf("Input: %s", s)

	var z bytes.Buffer
	sz := len(s)
	for i := 0; i < sz && i != -1; {
		r := strings.Index(s[i:], "src=\"/")
		if r >= 0 && r < sz {
			log.Printf("Found src at %d(%d)", r, i)
			z.WriteString(s[i : i+r])
			z.WriteString(fmt.Sprintf("src=\"%s/", url))
			if i+r+6 > len(s) {
				i = -1
			} else {
				i += r + 6
			}
		} else {
			z.WriteString(s[i:])
			i = -1
		}
	}

	s = z.String()
	//log.Printf("SRC Check: %s", s)

	var q bytes.Buffer
	sz = len(s)
	for i := 0; i < sz && i != -1; {
		r := strings.Index(s[i:], "href=\"/")
		if r >= 0 && r < sz {
			log.Printf("Found href at %d", r)
			q.WriteString(s[i : i+r])
			q.WriteString(fmt.Sprintf("href=\"%s/", url))
			if i+r+7 > len(s) {
				i = -1
			} else {
				i += r + 7
			}
		} else {
			q.WriteString(s[i:])
			i = -1
		}
	}

	log.Print(q.String())
	return q.String()
}
