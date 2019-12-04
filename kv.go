package xml

import (
	"bufio"
	"bytes"
)

// KV ...
type KV struct {
	k, v []byte
}

// Key ...
func (kv *KV) Key() string {
	return string(kv.k)
}

// KeyBytes ...
func (kv *KV) KeyBytes() []byte {
	return kv.k
}

// Value ...
func (kv *KV) Value() string {
	return string(kv.v)
}

// ValueBytes ...
func (kv *KV) ValueBytes() []byte {
	return kv.v
}

func (kv *KV) reset() {
	kv.k = kv.k[:0]
	kv.v = kv.v[:0]
}

func (kv *KV) parse(r *bufio.Reader) error {
	k, err := r.ReadBytes('=')
	if err == nil {
		kv.k = append(kv.k[:0], bytes.TrimRight(k[:len(k)-1], " \r\n")...)
		var (
			c byte
			v []byte
		)
	loop:
		for {
			c, err = skipWS(r)
			if err != nil {
				break
			}

			switch c {
			case '"':
				v, err = r.ReadBytes('"')
				if err == nil {
					kv.v = append(kv.v[:0], v[:len(v)-1]...)
				}
				break loop
			}
		}
	}
	return err
}
