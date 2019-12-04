package xml

import "bufio"

func skipWS(r *bufio.Reader) (c byte, err error) {
	for {
		c, err = r.ReadByte()
		if err != nil || c > 32 {
			break
		}
	}
	return
}
