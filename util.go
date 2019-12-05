package xml

import (
	"bufio"
	"unsafe"
)

func skipWS(r *bufio.Reader) (c byte, err error) {
	for {
		c, err = r.ReadByte()
		if err != nil || c > 32 {
			break
		}
	}
	return
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
