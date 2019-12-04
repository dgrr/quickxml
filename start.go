package xml

import (
	"bufio"
)

// StartElement ...
type StartElement struct {
	Name  string
	Attrs []*KV
}

func (s *StartElement) parse(r *bufio.Reader) error {
	c, err := skipWS(r) // skip any whitespaces
	if err != nil {
		return err
	}
	s.Name += string(c)
loop:
	for {
		c, err = r.ReadByte()
		if err != nil {
			break
		}
		switch c {
		case ' ', '>': // read until the first space or reaching the end
			break loop
		default:
			s.Name += string(c)
		}
	}
	if c == ' ' && err == nil { // doesn't reach the end
		err = s.parseAttrs(r)
	}

	return err
}

func (s *StartElement) parseAttrs(r *bufio.Reader) (err error) {
	var c byte
	for {
		c, err = skipWS(r) // skip whitespaces until reaching the key
		if err != nil || c == '>' {
			break
		}
		if c == '/' {
			continue
		}
		r.UnreadByte()

		// read key
		kv := new(KV)
		err = kv.parse(r)
		if err == nil {
			s.Attrs = append(s.Attrs, kv)
		}
	}
	return
}
