package xml

import (
	"bufio"
	"io"
)

// Reader represents a XML reader.
type Reader struct {
	r   *bufio.Reader
	err error
	e   Element
	n   *string
}

// NewReader returns a initialized reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// Element returns the last readed element.
func (r *Reader) Element() Element {
	return r.e
}

// Error return the last error.
func (r *Reader) Error() error {
	return r.err
}

func (r *Reader) release() {
	if r.e == nil {
		return
	}

	if e, ok := r.e.(*StartElement); ok {
		releaseStart(e)
	} else if e, ok := r.e.(*EndElement); ok {
		releaseEnd(e)
	}
	r.e = nil
}

// Next iterates until the next XML element.
func (r *Reader) Next() bool {
	r.release()

	var c byte
	for r.e == nil && r.err == nil {
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
	}

	return r.e != nil && r.err == nil
}

// AssignNext will assign the next TextElement to ptr.
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
