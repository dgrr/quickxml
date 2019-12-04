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

// ReleaseEnd ...
func ReleaseEnd(end *EndElement) {
	end.reset()
	endPool.Put(end)
}

// EndElement ...
type EndElement struct {
	name []byte
}

func (e *EndElement) reset() {
	e.name = e.name[:0]
}

// Name ...
func (e *EndElement) Name() string {
	return string(e.name)
}

// NameBytes ...
func (e *EndElement) NameBytes() []byte {
	return e.name
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
