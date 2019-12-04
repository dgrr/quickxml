package xml

import (
	"bufio"
	"strings"
)

// KV ...
type KV struct {
	K, V string
}

func (kv *KV) parse(r *bufio.Reader) (err error) {
	kv.K, err = r.ReadString('=')
	if err == nil {
		kv.K = strings.TrimRight(kv.K[:len(kv.K)-1], " \r\n")
		var c byte
	loop:
		for {
			c, err = skipWS(r)
			if err != nil {
				break
			}

			switch c {
			case '"':
				kv.V, err = r.ReadString('"')
				if err == nil {
					kv.V = strings.Trim(kv.V[:len(kv.V)-1], " ")
				}
				break loop
			}
		}
	}
	return
}
