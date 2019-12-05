package xml

import (
	"bufio"
	"sync"
)

var endPool = sync.Pool{
	New: func() interface{} {
		return new(EndElement)
	},
}

// ReleaseEnd returns an EndElement to the pool.
func ReleaseEnd(end *EndElement) {
	//end.reset()
	endPool.Put(end)
}

// EndElement represents a XML end element.
type EndElement struct {
	name []byte
}

func (e *EndElement) reset() {
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
