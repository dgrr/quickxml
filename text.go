package xml

import "bufio"

// TextElement represents a XML text.
type TextElement string

// NewText creates a new TextElement.
func NewText(str string) *TextElement {
	t := TextElement(str)
	return &t
}

func (t *TextElement) parse(_ *bufio.Reader) error {
	return nil
}

// String returns the string representation of TextElement.
func (t *TextElement) String() string {
	return string(*t)
}
