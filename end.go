package xml

import "bufio"

// EndElement ...
type EndElement struct {
	Name string
}

func (e *EndElement) parse(r *bufio.Reader) error {
	c, err := skipWS(r)
	if err != nil {
		return err
	}
	e.Name += string(c)
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
		e.Name += string(c)
	}

	return err
}
