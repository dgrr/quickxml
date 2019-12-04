package xml

import "bufio"

// TextElement ...
type TextElement string

func (t *TextElement) parse(_ *bufio.Reader) error {
	return nil
}
