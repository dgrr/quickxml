package xml

import "bufio"

// Element represents a XML element.
//
// Element can be:
// - StartElement.
// - EndElement.
// - TextElement.
type Element interface {
	parse(r *bufio.Reader) error
	String() string
}
