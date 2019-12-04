package main

import (
	"fmt"
	"strings"

	"github.com/dgrr/xml"
)

// Book represents our XML structure
type Book struct {
	Category string
	Title    string
	Author   string
	Year     string
	Price    string
}

// String will print the book info in a string (for fmt)
func (book Book) String() string {
	return fmt.Sprintf(
		"%s\n  Title: %s\n  Author: %s\n  Year: %s\n  Price: %s",
		book.Category, book.Title, book.Author, book.Year, book.Price,
	)
}

func main() {
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

	book := Book{}
	r := xml.NewReader(strings.NewReader(str))
	for r.Next() {
		switch e := r.Element().(type) {
		case *xml.StartElement:
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
		case *xml.EndElement:
			if e.Name() == "book" { // book parsed
				fmt.Printf("%s\n", book)
			}
		}
	}
}
