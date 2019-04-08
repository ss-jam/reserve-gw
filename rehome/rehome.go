package rehome

import (
	"bytes"
	"fmt"
	"log"
)

// fix any src or href attributes that are relative to the reference URL
func FixAttributes(b []byte, url string) string {
	//	h := strings.Index(s, "<head")
	//	he := strings.Index(s, "</head>")
	//	b := strings.Index(s, "<body")
	//	be := strings.Index(s, "</body>")
	//s := string(b[:])
	//log.Printf("Input: %s", s)

	var s bytes.Buffer
	sz := len(b)
	for i := 0; i < sz && i != -1; {
		r := bytes.Index(b[i:], []byte("src=\"/"))
		if r >= 0 && r < sz {
			log.Printf("Found src at %d(%d)", r, i)
			s.Write(b[i : i+r])
			s.Write([]byte(fmt.Sprintf("src=\"%s/", url)))
			if i+r+6 > len(b) {
				i = -1
			} else {
				i += r + 6
			}
		} else {
			s.Write(b[i:])
			i = -1
		}
	}

	//s = z.String()
	//log.Printf("SRC Check: %s", s)

	var q bytes.Buffer
	sz = s.Len()
	for i := 0; i < sz && i != -1; {
		r := bytes.Index(s.Bytes()[i:], []byte("href=\"/"))
		if r >= 0 && r < sz {
			log.Printf("Found href at %d", r)
			q.Write(s.Bytes()[i : i+r])
			q.Write([]byte(fmt.Sprintf("href=\"%s/", url)))
			if i+r+7 > s.Len() {
				i = -1
			} else {
				i += r + 7
			}
		} else {
			q.Write(s.Bytes()[i:])
			i = -1
		}
	}

	log.Print(q.String())
	return q.String()
}
