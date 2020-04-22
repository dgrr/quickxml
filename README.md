# Quick XML

[![Go Report Card](https://goreportcard.com/badge/github.com/dgrr/xml)](https://goreportcard.com/report/github.com/dgrr/xml)
[![Build Status](https://travis-ci.com/dgrr/xml.svg?branch=master)](https://travis-ci.com/dgrr/xml)
[![codecov](https://codecov.io/gh/dgrr/xml/branch/master/graph/badge.svg)](https://codecov.io/gh/dgrr/xml)

QuickXML is a package to process XML files in an iterative way. It doesn't use reflect so you'll need to work a little more :D

Most of the times working with XML is a painful task. Also, the Golang std library doesn't help too much. Neither is fast nor has good doc. This library just tries to process XML files in a iterative way, ignoring most of the common errors in XML (like non-closed tags or optional tags). So it just detects when a tag is being open and closed, and doesn't have control whether the tag X has been open before closed or viceversa.

**IMPORTANT NOTE: This package doesn't provide a fully featured XML. It has been created for XLSX parsing.**

PRs are welcome.

# How QuickXML performs?


![alt_text](https://github.com/dgrr/quickxml/blob/master/_imgs/alloc_time.png)

The graph shows the amount of memory used (blue) represented in the left axis.
And the time spent in seconds (red) represented in the right axis.

As you can see, QuickXML is the fastest and the second one on memory usage.in

# Example


```go
package main

import (
	"fmt"
	"strings"

	xml "github.com/dgrr/quickxml"
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
```
