// +build ignore
package main

import (
	"os"

	"github.com/dgrr/xml"
)

func main() {
	w := xml.NewWriter(os.Stdout)
	es := []xml.Element{
		xml.NewStart("s", true, nil),
		xml.NewStart("ss", false, xml.NewAttrs("a1", "a2", "a3", "a4")),
		xml.NewText("This is a middle next"),
		xml.NewStart("sss", false, xml.NewAttrs("k", "v", "kk")),
		xml.NewText("Another text"),
		xml.NewEnd("sss"),
		xml.NewEnd("ss"),
	}

	for _, e := range es {
		w.WriteIndent(e)
	}
}
