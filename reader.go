package xml

import (
	"bufio"
	"io"
)

// Element ...
type Element interface {
	parse(r *bufio.Reader) error
}

// Reader ...
type Reader struct {
	r   *bufio.Reader
	err error
	e   Element
	n   *string
}

// NewReader ...
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// Element ...
func (r *Reader) Element() Element {
	return r.e
}

// Err ...
func (r *Reader) Err() error {
	return r.err
}

// Next ...
func (r *Reader) Next() bool {
	var c byte
	r.e = nil
	for r.e == nil {
		c, r.err = skipWS(r.r)
		if r.err == nil {
			switch c { // get next token
			case '<': // new element
				r.next()
			default: // text string
				r.r.UnreadByte()
				t, err := r.r.ReadString('<') // read until a new element starts (or EOF is reached)
				if err != nil {
					r.err = err
				} else {
					t = t[:len(t)-1]
					if r.n != nil {
						*r.n, r.n = t, nil
					} else {
						tt := TextElement(t)
						r.e = &tt
					}
					r.r.UnreadByte()
				}
			}
		}
		if r.err != nil {
			break
		}
	}

	return r.e != nil
}

// AssignNext ...
func (r *Reader) AssignNext(ptr *string) {
	r.n = ptr
}

// skip reads until the next end tag '>'
func (r *Reader) skip() error {
	_, err := r.r.ReadBytes('>')
	return err
}

// next will read the next byte after finding '<'
func (r *Reader) next() {
	var c byte
	c, r.err = skipWS(r.r)
	if r.err == nil {
		switch c {
		case '/':
			r.e = endPool.Get().(*EndElement)
		case '!':
			r.err = r.skip()
		case '?':
			r.err = r.skip()
		default:
			r.e = startPool.Get().(*StartElement)
			r.r.UnreadByte()
		}
		if r.err == nil && r.e != nil {
			r.err = r.e.parse(r.r)
			if r.err != nil {
				r.e = nil
			}
		}
	}
}
