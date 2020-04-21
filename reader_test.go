package xml

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
	"testing"
)

func TestReaderBasic(t *testing.T) {
	const str = `<?xml version="1.0" encoding="UTF-8"?>
	<   first  >
		<  second     k  =  "v"    k2  =  "v2"   >text<  /  second  >
	<  /  first  >`

	start := []string{"first", "second"}
	text := map[int]string{
		2: "text",
	}

	attrs := map[int][]KV{
		1: {
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
				e.Attrs().RangeWithIndex(func(i int, kv *KV) {
					if !bytes.Equal(ekv[i].KeyBytes(), kv.KeyBytes()) {
						t.Fatalf("Unexpected Attr Key on %d: got `%s`. Expected `%s`. Len %d", i, kv.Key(), ekv[i].Key(), e.Attrs().Len())
					}
					if !bytes.Equal(ekv[i].ValueBytes(), kv.ValueBytes()) {
						t.Fatalf("Unexpected Attr Value on %d: got `%s`. Expected `%s`", i, kv.Value(), ekv[i].Value())
					}
				})
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

func TestXMLContentType(t *testing.T) {
	const xmlStr = `<?xml version="1.0" encoding="utf-8" standalone="yes"?>

  <Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
    <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
    <Override PartName="/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
    <Override PartName="/xl/worksheets/data.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
    <Override PartName="/stylesheet.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/></Types>`

	lookFor := `application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml`

	r := NewReader(strings.NewReader(xmlStr))
	for r.Next() {
		if se, ok := r.Element().(*StartElement); ok {
			if kv := se.Attrs().Get("ContentType"); kv != nil && kv.Value() == lookFor {
				return
			}
		}
	}

	t.Fatalf("%s not found", lookFor)
}

// benchmark
// Book represents our XML structure
type Book struct {
	XMLName  xml.Name `xml:"book"`
	Category string   `xml:"category,attr"`
	Title    string   `xml:"title"`
	Author   string   `xml:"author"`
	Year     string   `xml:"year"`
	Price    string   `xml:"price"`
}

const benchStr = `<bookstore xmlns:p="urn:schemas-books-com:prices">

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

type Buffer struct {
	b []byte
	i int
}

func (bf *Buffer) Read(b []byte) (int, error) {
	if len(bf.b) == bf.i {
		return 0, io.EOF
	}

	n := copy(b, bf.b[bf.i:])
	bf.i += n
	return n, nil
}

func (bf *Buffer) Reset() {
	bf.i = 0
}

func BenchmarkFastXML(b *testing.B) {
	sr := &Buffer{
		b: []byte(benchStr),
		i: 0,
	}
	for i := 0; i < b.N; i++ {
		r := NewReader(sr)
		benchFastXML(b, r)
		sr.Reset()
	}
}

func benchFastXML(b *testing.B, r *Reader) {
	books := 0
	book := Book{}
	for r.Next() {
		switch e := r.Element().(type) {
		case *StartElement:
			switch e.Name() {
			case "book": // start reading a book
				attr := e.Attrs().Get("category")
				if attr != nil {
					book.Category = attr.Value()
				}
			case "title": // You can capture the lang too using e.Attrs()
				r.AssignNext(&book.Title)
				// AssignNext will assign the next text found to book.Title
			case "author":
				r.AssignNext(&book.Author)
			case "year":
				r.AssignNext(&book.Year)
			case "p:price":
				r.AssignNext(&book.Price)
			}
			ReleaseStart(e)
		case *EndElement:
			if e.Name() == "book" { // book parsed
				books++
			}
			ReleaseEnd(e)
		}
	}
	if r.Error() != nil && r.Error() != io.EOF {
		b.Fatal(r.Error())
	}
	if books != 4 {
		b.Fatalf("Expected 4 books. Got %d", books)
	}
}

func BenchmarkXML(b *testing.B) {
	sr := &Buffer{
		b: []byte(benchStr),
		i: 0,
	}
	for i := 0; i < b.N; i++ {
		d := xml.NewDecoder(sr)
		benchXML(b, d)
		sr.Reset()
	}
}

func benchXML(b *testing.B, d *xml.Decoder) {
	books := 0
	book := Book{}
	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			b.Fatal(err)
		}

		switch e := tok.(type) {
		case xml.StartElement:
			switch e.Name.Local {
			case "book":
				for _, attr := range e.Attr {
					if attr.Name.Local == "category" {
						book.Category = attr.Value
					}
				}
			case "title":
				d.DecodeElement(&book.Title, &e)
			case "author":
				d.DecodeElement(&book.Author, &e)
			case "year":
				d.DecodeElement(&book.Year, &e)
			case "price":
				d.DecodeElement(&book.Price, &e)
			}
		case xml.EndElement:
			if e.Name.Local == "book" {
				books++
			}
		}
	}
	if books != 4 {
		b.Fatalf("Expected 4 books. Got %d", books)
	}
}
