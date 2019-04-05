// Manage the HTML elements
package elements

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type Pdoc struct {
	Head []byte
	Body []byte
}

// Process an HTML string
func BatchProcess(s []byte, root string) (*Pdoc, error) {
	z := html.NewTokenizer(strings.NewReader(string(s)))

	h := make([]byte, 1024*100)
	b := make([]byte, 1024*100)
	done := false

	inHead := false
	inBody := false
	hindex := 0
	bindex := 0

	for !done {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return nil, z.Err()
		case html.TextToken:
			// emitBytes should copy the []byte it receives,
			// if it doesn't process it immediately.
			log.Printf("Text: %s", z.Text())
			raw := z.Raw()
			if inHead {
				copy(h[hindex:], raw)
				hindex = hindex + len(raw)
			} else if inBody {
				copy(b[bindex:], raw)
				bindex = bindex + len(raw)
			}
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			switch string(tn) {
			case "html":
				if tt == html.EndTagToken {
					done = true
				}
			case "head":
				log.Printf("HEAD: %s", z.Raw())
				if tt == html.StartTagToken {
					inHead = true
					//h = z.Raw()
				} else {
					inHead = false
				}
			case "body":
				log.Printf("BODY: %s", z.Raw())
				if tt == html.StartTagToken {
					inBody = true
					//b = z.Raw()
				} else {
					inBody = false
				}
			default:
				raw := z.Raw()
				log.Printf("%s: %s", strings.ToUpper(string(tn)), raw)
				if inHead {
					copy(h[hindex:], raw)
					hindex = hindex + len(raw)
				} else if inBody {
					copy(b[bindex:], raw)
					bindex = bindex + len(raw)
				}
			}
		}
	}
	//log.Printf("Pdoc.Head -> %s", h)
	return &Pdoc{h[:hindex], b[:bindex]}, nil
}

func emitNode(n *html.Node) string {
	s := ""
	if len(n.Attr) > 0 {
		for _, a := range n.Attr {
			s = s + fmt.Sprintf(" %s=%s,", a.Key, a.Val)
		}
		s = s[:len(s)-1]
	}
	var r string
	switch n.Type {
	case html.TextNode:
		r = n.Data
	case html.ElementNode:
		r = fmt.Sprintf("<%s%s>", n.DataAtom, s)
	case html.DocumentNode:
		log.Printf("Document Node: %s, %s, %v", n.DataAtom, n.Data, n.Attr)
		r = n.Data
	case html.CommentNode:
		log.Printf("Commnet Node: %s, %s, %v", n.DataAtom, n.Data, n.Attr)
		r = fmt.Sprintf("<!%s", n.Data)
	case html.DoctypeNode:
		log.Printf("Doctype Node: %s, %s, %v", n.DataAtom, n.Data, n.Attr)
		r = n.Data
	}
	return r
}

type Mdoc map[string]string

func logger(f string, da atom.Atom) {
	ignore := map[atom.Atom]string{
		atom.Script: "script",
		atom.Li:     "li",
		atom.A:      "a",
		atom.Div:    "div",
		atom.P:      "p",
	}

	_, found := ignore[da]
	if !found && len(da.String()) > 0 {
		log.Print(f, da)
	}
}

// Parse an HTML slice
func BatchParse(s []byte, root string) (*Pdoc, error) {
	d := make([]byte, 1024*100)
	index := 0
	bindex := 0

	doc, err := html.Parse(strings.NewReader(string(s)))
	if err != nil {
		log.Printf("Parse error: %s", err)
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		logger("Open ", n.DataAtom)
		if n.Type == html.ErrorNode {
			log.Printf("Unfortunate error node: '%s', '%s'", n.Data, n.DataAtom)
		} else {
			logger(fmt.Sprintf("%s: %s %s (%d)", n.Data, n.Attr, n.DataAtom, n.Type), n.DataAtom)
			if n.DataAtom == atom.Body {
				bindex = index
			} else if n.DataAtom == atom.Head {
				index = 0
			} else if n.DataAtom != atom.Html {
				l := emitNode(n)
				copy(d[index:], l)
				index += len(l)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if c.NextSibling == nil && len(n.DataAtom.String()) > 0 {
				switch n.DataAtom {
				case atom.Img, atom.Image, atom.Br:
					log.Printf("Close %s but do not write", n.DataAtom)
				default:
					logger("Close ", n.DataAtom)
					copy(d[index:], fmt.Sprintf("</%s>", n.DataAtom))
					index += len(n.DataAtom.String()) + 3
				}
			}
		}
	}
	f(doc)
	return &Pdoc{d[:bindex], d[bindex:index]}, nil
}

func StreamProcess(resp *http.Response) *Pdoc {
	pdoc := Pdoc{}
	doc := html.NewTokenizer(resp.Body)
	for tokenType := doc.Next(); tokenType != html.ErrorToken; {
		token := doc.Token()
		if tokenType == html.StartTagToken {
			if token.DataAtom != atom.A {
				tokenType = doc.Next()
				continue
			}
			//Do something with it, for example extract url
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					//url here attr.Val
					//ideally send it to some worker instance to avoid blocking here
				}
			}
		}
		tokenType = doc.Next()
	}
	resp.Body.Close()
	return &pdoc
}
