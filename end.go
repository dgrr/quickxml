package xml

import (
	"bufio"
	"fmt"
	"sync"
)

var endPool = sync.Pool{
	New: func() interface{} {
		return new(EndElement)
	},
}

// ReleaseEnd returns an EndElement to the pool.
func ReleaseEnd(end *EndElement) {
	end.Reset()
	endPool.Put(end)
}

// EndElement represents a XML end element.
type EndElement struct {
	name []byte
}

// NewEnd creates a new EndElement.
func NewEnd(name string) *EndElement {
	return &EndElement{
		name: []byte(name),
	}
}

// String returns the string representation of EndElement.
func (e *EndElement) String() string {
	return fmt.Sprintf("</%s>", e.name)
}

// SetName sets the name to the end element.
func (e *EndElement) SetName(name string) {
	e.name = []byte(name)
}

// SetNameBytes sets the name to the end element in bytes.
func (e *EndElement) SetNameBytes(name []byte) {
	e.name = append(e.name[:0], name...)
}

func (e *EndElement) Reset() {
	e.name = e.name[:0]
}

// Name returns the name of the XML node.
func (e *EndElement) Name() string {
	return string(e.name)
}

// NameBytes returns the name of the XML node in bytes.
func (e *EndElement) NameBytes() []byte {
	return e.name
}

// NameUnsafe returns a string holding the name parameter.
//
// This function differs from Name() on using unsafe methods.
func (e *EndElement) NameUnsafe() string {
	return b2s(e.name)
}

func (e *EndElement) parse(r *bufio.Reader) error {
	c, err := skipWS(r)
	if err != nil {
		return err
	}
	e.name = append(e.name[:0], c)
	for {
		c, err = r.ReadByte()
		if err != nil {
			break
		}
		if c == '>' {
			break
		}
		if c == ' ' {
			continue
		}
		e.name = append(e.name, c)
	}

	return err
}
