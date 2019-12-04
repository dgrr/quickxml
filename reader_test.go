package xml

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestReaderBasic(t *testing.T) {
	const str = `<?xml version="1.0" encoding="UTF-8"?>
	<first>
		<second k="v" k2="v2">text</second>
	</first>`

	start := []string{"first", "second"}
	text := map[int]string{
		2: "text",
	}

	attrs := map[int][]KV{
		1: []KV{
			{
				k: []byte("k"),
				v: []byte("v"),
			},
			{
				k: []byte("k2"),
				v: []byte("v2"),
			},
		},
	}

	starti := 0

	r := NewReader(strings.NewReader(str))
	for r.Next() {
		switch e := r.Element().(type) {
		case *StartElement:
			if e.Name() != start[starti] {
				t.Fatalf("Unexpected StartElement: got `%s`. Expected `%s`", e.Name(), start[starti])
			}
			ekv, ok := attrs[starti]
			if ok {
				for i, kv := range e.Attrs() {
					if !bytes.Equal(ekv[i].KeyBytes(), kv.KeyBytes()) {
						t.Fatalf("Unexpected Attr Key on %d: got `%s`. Expected `%s`. Len %d", i, kv.Key(), ekv[i].Key(), len(e.Attrs()))
					}
					if !bytes.Equal(ekv[i].ValueBytes(), kv.ValueBytes()) {
						t.Fatalf("Unexpected Attr Value on %d: got `%s`. Expected `%s`", i, kv.Value(), ekv[i].Value())
					}
				}
			}
			ReleaseStart(e)
			starti++
		case *TextElement:
			s, ok := text[starti]
			if !ok {
				t.Fatalf("Expected `%s` on %d. Got `%s`", s, starti, *e)
			} else if s != string(*e) {
				t.Fatalf("Unexpected text. Got `%s`. Expected `%s`", *e, s)
			}
		case *EndElement:
			starti--
			if e.Name() != start[starti] {
				t.Fatalf("Unexpected EndElement: got `%s`. Expected `%s`", e.Name(), start[starti])
			}
			ReleaseEnd(e)
		}
	}
}

func TestReaderComplex(t *testing.T) {
	const str = `<bookstore xmlns:p="urn:schemas-books-com:prices">

	<book category="COOKING">
	  <title lang="en">Everyday Italian</title>
	  <author>Giada De Laurentiis</author>
	  <year>2005</year>
	  <p:price>30.00</p:price>
	</book>
  
	<book category="CHILDREN">
	  <title lang="en">Harry Potter</title>
	  <author>J K. Rowling</author>
	  <year>2005</year>
	  <p:price>29.99</p:price>
	</book>
  
	<book category="WEB">
	  <title lang="en">XQuery Kick Start</title>
	  <author>James McGovern</author>
	  <author>Per Bothner</author>
	  <author>Kurt Cagle</author>
	  <author>James Linn</author>
	  <author>Vaidyanathan Nagarajan</author>
	  <year>2003</year>
	  <p:price>49.99</p:price>
	</book>
  
	<book category="WEB">
	  <title lang="en">Learning XML</title>
	  <author>Erik T. Ray</author>
	  <year>2003</year>
	  <p:price>39.95</p:price>
	</book>
  
  </bookstore>`

	r := NewReader(strings.NewReader(str))
	spaces := ""
	for r.Next() {
		switch e := r.Element().(type) {
		case *StartElement:
			fmt.Printf("%s<%s", spaces, e.Name())
			for _, kv := range e.Attrs() {
				fmt.Printf(` %s="%s"`, kv.Key(), kv.Value())
			}
			fmt.Println(">")
			spaces += "  "
			ReleaseStart(e)
		case *TextElement:
			fmt.Printf("%s%s\n", spaces, *e)
		case *EndElement:
			spaces = spaces[:len(spaces)-2]
			fmt.Printf("%s</%s>\n", spaces, e.Name())
			ReleaseEnd(e)
		}
	}
}
