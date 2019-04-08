// Simple HTML tree parser for manipulating a page
package manipulate

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type MyNode struct {
	DataAtom string
	Raw      []byte
}

func fixUrl(raw []byte, attrs []html.Attribute, look4 string, url string) []byte {
	var buf bytes.Buffer
	found := false
	for _, attr := range attrs {
		if attr.Key == look4 && attr.Val[0] == '/' {
			i := bytes.Index(raw, []byte(fmt.Sprintf("%s=\"", look4)))
			buf.Write(raw[:i])
			buf.Write([]byte(fmt.Sprintf("%s=\"%s", look4, url)))
			buf.Write(raw[i+len(look4)+2:])
			found = true
		}
	}
	if !found {
		return raw
	} else {
		return buf.Bytes()
	}
}

func SimpleBatch(b []byte, url string, ref string) string {
	z := html.NewTokenizer(bytes.NewReader(b))

	done := false
	var buf bytes.Buffer
	for !done {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				done = true
			} else {
				log.Printf("Error on node: %s", z.Err())
			}
		// The default case, but expanded here for clarity
		case html.TextToken, html.DoctypeToken, html.CommentToken, html.EndTagToken:
			raw := z.Raw()
			buf.Write(raw)
		case html.SelfClosingTagToken:
			raw := z.Raw()
			t := z.Token()
			log.Printf("Self tag: %s: %s", t.DataAtom, raw)
			switch t.DataAtom {
			case atom.Img:
				buf.Write(fixUrl(raw, t.Attr, "src", url))
			default:
				buf.Write(raw)
			}
		case html.StartTagToken:
			raw := z.Raw()
			t := z.Token()
			log.Printf("Start tag: %s: %s", t.DataAtom, raw)
			switch t.DataAtom {
			case atom.Script:
				buf.Write(fixUrl(raw, t.Attr, "src", url))
			// Link nodes in HEAD are not self closing even though you may think they should be
			case atom.Link:
				buf.Write(fixUrl(raw, t.Attr, "href", url))
			case atom.A:
				buf.Write(fixUrl(raw, t.Attr, "href", ref))
			case atom.Img:
				buf.Write(fixUrl(raw, t.Attr, "src", url))
			default:
				buf.Write(raw)
			}
		}
	}
	return buf.String()
}
