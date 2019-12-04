package xml

import (
	"bytes"
	"strings"
	"testing"
)

func TestReaderBasic(t *testing.T) {
	const str = `<?xml version="1.0" encoding="UTF-8"?>
	<first>
		<second k="v" k2="v2">text</second>
	</first>`

	start := []string{"first", "second"}
	text := map[int]string{
		2: "text",
	}

	attrs := map[int][]KV{
		1: []KV{
			{
				k: []byte("k"),
				v: []byte("v"),
			},
			{
				k: []byte("k2"),
				v: []byte("v2"),
			},
		},
	}

	starti := 0

	r := NewReader(strings.NewReader(str))
	for r.Next() {
		switch e := r.Element().(type) {
		case *StartElement:
			if e.Name() != start[starti] {
				t.Fatalf("Unexpected StartElement: got `%s`. Expected `%s`", e.Name(), start[starti])
			}
			ekv, ok := attrs[starti]
			if ok {
				e.Attrs().RangeWithIndex(func(i int, kv *KV) {
					if !bytes.Equal(ekv[i].KeyBytes(), kv.KeyBytes()) {
						t.Fatalf("Unexpected Attr Key on %d: got `%s`. Expected `%s`. Len %d", i, kv.Key(), ekv[i].Key(), e.Attrs().Len())
					}
					if !bytes.Equal(ekv[i].ValueBytes(), kv.ValueBytes()) {
						t.Fatalf("Unexpected Attr Value on %d: got `%s`. Expected `%s`", i, kv.Value(), ekv[i].Value())
					}
				})
			}
			ReleaseStart(e)
			starti++
		case *TextElement:
			s, ok := text[starti]
			if !ok {
				t.Fatalf("Expected `%s` on %d. Got `%s`", s, starti, *e)
			} else if s != string(*e) {
				t.Fatalf("Unexpected text. Got `%s`. Expected `%s`", *e, s)
			}
		case *EndElement:
			starti--
			if e.Name() != start[starti] {
				t.Fatalf("Unexpected EndElement: got `%s`. Expected `%s`", e.Name(), start[starti])
			}
			ReleaseEnd(e)
		}
	}
}
