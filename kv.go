package xml

import (
	"bufio"
)

// KV represents an attr which is a key-value pair.
type KV struct {
	k, v []byte
}

// Key returns the key.
func (kv *KV) Key() string {
	return string(kv.k)
}

// KeyBytes returns the key.
func (kv *KV) KeyBytes() []byte {
	return kv.k
}

// KeyUnsafe returns a string holding the name parameter.
//
// This function differs from Key() on using unsafe methods.
func (kv *KV) KeyUnsafe() string {
	return b2s(kv.k)
}

// Value returns the value.
func (kv *KV) Value() string {
	return string(kv.v)
}

// ValueBytes returns the value.
func (kv *KV) ValueBytes() []byte {
	return kv.v
}

// ValueUnsafe returns a string holding the name parameter.
//
// This function differs from Value() on using unsafe methods.
func (kv *KV) ValueUnsafe() string {
	return b2s(kv.v)
}

func (kv *KV) reset() {
	kv.k = kv.k[:0]
	kv.v = kv.v[:0]
}

func (kv *KV) parse(r *bufio.Reader) error {
	k, err := r.ReadBytes('=')
	if err == nil {
		n := len(k) - 2
		for k[n] == ' ' {
			n--
		}

		kv.k = append(kv.k[:0], k[:n+1]...)
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
