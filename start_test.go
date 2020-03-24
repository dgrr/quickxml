package xml

import (
	"strings"
	"testing"
)

func checkXML(t *testing.T, xmlStr, eName string, attMap map[string]string, hasEnd bool) {
	r := NewReader(strings.NewReader(xmlStr))
	if !r.Next() {
		t.Fatal("Next() == false")
	}

	e := r.Element()
	if se, ok := e.(*StartElement); !ok {
		t.Fatal("Element() != *StartElement")
	} else {
		if se.Name() != eName {
			t.Fatalf("%s != %s", se.Name(), eName)
		}

		if attMap != nil {
			att := se.Attrs()
			att.Range(func(kv *KV) {
				v, ok := attMap[kv.Key()]
				if !ok {
					t.Fatalf("%s doesn't match any StartElement attr", kv.Key())
				}
				if v != kv.Value() {
					t.Fatalf("%s != %s", kv.Value(), v)
				}
			})
		}

		if se.HasEnd() != hasEnd {
			t.Fatalf("StartElement end expected %v. Got %v", hasEnd, se.HasEnd())
		}
	}
}

func TestStartElement(t *testing.T) {
	const xmlStr = `  <   element  >`

	checkXML(t, xmlStr, "element", nil, false)
}

func TestStartElementWithEnd(t *testing.T) {
	const xmlStr = `  <   element />`

	checkXML(t, xmlStr, "element", nil, true)
}

func TestStartElementWithAttrs(t *testing.T) {
	const xmlStr = `  <   element attr1="a" attr2   =  "b"  >`
	attMap := map[string]string{
		"attr1": "a",
		"attr2": "b",
	}

	checkXML(t, xmlStr, "element", attMap, false)
}

func TestStartElementWithAttrsWithEnd(t *testing.T) {
	const xmlStr = `  <   element attr1="a" attr2   =  "b"  />`
	attMap := map[string]string{
		"attr1": "a",
		"attr2": "b",
	}

	checkXML(t, xmlStr, "element", attMap, true)
}
