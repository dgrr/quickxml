package xml

import (
	"bufio"
	"sync"
)

var startPool = sync.Pool{
	New: func() interface{} {
		return new(StartElement)
	},
}

// ReleaseStart ...
func ReleaseStart(start *StartElement) {
	//start.reset()
	startPool.Put(start)
}

// StartElement ...
type StartElement struct {
	name  []byte
	attrs []KV
}

// Name ...
func (s *StartElement) Name() string {
	return string(s.name)
}

// NameBytes ...
func (s *StartElement) NameBytes() []byte {
	return s.name
}

// Attrs ...
func (s *StartElement) Attrs() []KV {
	return s.attrs
}

func (s *StartElement) reset() {
	s.name = s.name[:0]
	s.attrs = s.attrs[:0]
}

func (s *StartElement) parse(r *bufio.Reader) error {
	c, err := skipWS(r) // skip any whitespaces
	if err != nil {
		return err
	}
	s.name = append(s.name[:0], c)
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
			s.name = append(s.name, c)
		}
	}
	if c == ' ' && err == nil { // doesn't reach the end
		s.attrs = s.attrs[:0]
		err = s.parseAttrs(r)
	}

	return err
}

func (s *StartElement) parseAttrs(r *bufio.Reader) (err error) {
	var c byte
	idx := 0
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
		err = s.getNextElement(idx).parse(r)
		if err == nil {
			idx++
		}
	}
	return
}

func (s *StartElement) getNextElement(idx int) *KV {
	if idx < cap(s.attrs) {
		s.attrs = s.attrs[:idx+1]
	} else {
		s.attrs = append(s.attrs, make([]KV, idx+1-cap(s.attrs))...)
	}

	return &s.attrs[idx]
}
