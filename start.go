package xml

import (
	"bufio"
	"bytes"
	"fmt"
	"sync"
)

var startPool = sync.Pool{
	New: func() interface{} {
		return new(StartElement)
	},
}

// releaseStart returns the StartElement to the pool.
func releaseStart(start *StartElement) {
	//start.reset()
	startPool.Put(start)
}

// Attrs represents the attributes of an XML StartElement.
type Attrs []KV

// NewAttr creates a new attribyte list.
//
// The attributes are a key-value pair of strings.
// For example: NewAttrs("k", "v") will create the attr k="v".
//
// If the attrs are odd nothing happens. The value associated
// with that key will be empty.
func NewAttrs(attrs ...string) *Attrs {
	att := make(Attrs, (len(attrs)/2)+(1&len(attrs)))

	for i, attr := range attrs {
		if i&1 == 0 {
			att[i/2].k = append(att[i/2].k, attr...)
		} else {
			att[i/2].v = append(att[i/2].v, attr...)
		}
	}

	return &att
}

// CopyTo copies kvs to kv2.
func (kvs *Attrs) CopyTo(kv2 *Attrs) {
	if n := len(*kvs) - len(*kv2); n > 0 {
		*kv2 = append(*kv2, make([]KV, n)...)
	}

	kvs.RangeWithIndex(func(i int, kv *KV) {
		(*kv2)[i].k = append((*kv2)[i].k[:0], kv.k...)
		(*kv2)[i].v = append((*kv2)[i].v[:0], kv.v...)
	})
}

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

// NewStart creats a new StartElement.
func NewStart(name string, hasEnd bool, attrs *Attrs) *StartElement {
	s := &StartElement{
		name:   []byte(name),
		hasEnd: hasEnd,
	}
	if attrs != nil {
		attrs.CopyTo(&s.attrs)
	}

	return s
}

func (s *StartElement) String() string {
	str := fmt.Sprintf("<%s", s.name)
	s.attrs.Range(func(kv *KV) {
		str += fmt.Sprintf(` %s="%s"`, kv.Key(), kv.Value())
	})
	if s.hasEnd {
		str += "/>"
	} else {
		str += ">"
	}
	return str
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

// SetName sets a string as StartElement's name.
func (s *StartElement) SetName(name string) {
	s.name = []byte(name)
}

// SetNameBytes sets the name bytes to the StartElement.
func (s *StartElement) SetNameBytes(name []byte) {
	s.name = append(s.name[:0], name...)
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

// Reset sets the default values to the StartElement.
func (s *StartElement) Reset() {
	s.name = s.name[:0]
	s.attrs = s.attrs[:0]
	s.hasEnd = false
}

func (s *StartElement) parse(r *bufio.Reader) error {
	s.Reset()

	c, err := skipWS(r) // skip any whitespaces
	if err != nil {
		return err
	}
	s.name = append(s.name, c)

	for {
		c, err = r.ReadByte()
		if err != nil || c == ' ' || c == '>' {
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
