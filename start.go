package xml

import (
	"bufio"
	"bytes"
	"sync"
)

var startPool = sync.Pool{
	New: func() interface{} {
		return new(StartElement)
	},
}

// ReleaseStart returns the StartElement to the pool.
func ReleaseStart(start *StartElement) {
	//start.reset()
	startPool.Put(start)
}

// Attrs represents the attributes of an XML StartElement.
type Attrs []KV

// Len returns the number of attributes.
func (kvs *Attrs) Len() int {
	return len(*kvs)
}

// Get returns the attribute based on name.
//
// If the name doesn't match any of the keys KV will be nil.
func (kvs *Attrs) Get(name string) *KV {
	for _, kv := range *kvs {
		if kv.KeyUnsafe() == name {
			return &kv
		}
	}
	return nil
}

// GetBytes returns the attribute based on name.
//
// If the name doesn't match any of the keys KV will be nil.
func (kvs *Attrs) GetBytes(name []byte) *KV {
	for _, kv := range *kvs {
		if bytes.Equal(kv.KeyBytes(), name) {
			return &kv
		}
	}
	return nil
}

// Range passes every attr to fn.
func (kvs *Attrs) Range(fn func(kv *KV)) {
	for _, kv := range *kvs {
		fn(&kv)
	}
}

// RangePre passes every attr to fn.
//
// If fn returns false the range loop will break.
func (kvs *Attrs) RangePre(fn func(kv *KV) bool) {
	for _, kv := range *kvs {
		if !fn(&kv) {
			break
		}
	}
}

// RangeWithIndex passes every attr and the index to fn.
func (kvs *Attrs) RangeWithIndex(fn func(i int, kv *KV)) {
	for i, kv := range *kvs {
		fn(i, &kv)
	}
}

// StartElement represents the start of a XML node.
type StartElement struct {
	name   []byte
	attrs  Attrs
	hasEnd bool
}

// HasEnd indicates if the StartElement ends as />
// Having this true means we do not expect a EndElement.
func (s *StartElement) HasEnd() bool {
	return s.hasEnd
}

// Name returns the name of the element.
func (s *StartElement) Name() string {
	return string(s.name)
}

// NameBytes returns the name of the element.
func (s *StartElement) NameBytes() []byte {
	return s.name
}

// NameUnsafe returns a string holding the name parameter.
//
// This function differs from Name() on using unsafe methods.
func (s *StartElement) NameUnsafe() string {
	return b2s(s.name)
}

// Attrs returns the attributes of an element.
func (s *StartElement) Attrs() *Attrs {
	return &s.attrs
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

	for {
		c, err = r.ReadByte()
		if err != nil {
			break
		}
		if c == ' ' || c == '>' {
			break
		}

		switch c {
		case '/':
			s.hasEnd = true
		default:
			if s.hasEnd { // malformed ??
				continue
			}
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
			s.hasEnd = true
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
