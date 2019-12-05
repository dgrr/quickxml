package xml

import "bufio"

// TextElement represents a XML text.
type TextElement string

func (t *TextElement) parse(_ *bufio.Reader) error {
	return nil
}
