// +build ignore
package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
)

type location struct {
	Data string `xml:",chardata"`
}

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	PrintMemUsage()

	d := xml.NewDecoder(f)
	count := 0
	for {
		tok, err := d.Token()
		if tok == nil || err == io.EOF {
			// EOF means we're done.
			break
		} else if err != nil {
			log.Fatalf("Error decoding token: %s", err)
		}

		switch ty := tok.(type) {
		case xml.StartElement:
			if ty.Name.Local == "location" {
				// If this is a start element named "location", parse this element
				// fully.
				var loc location
				if err = d.DecodeElement(&loc, &ty); err != nil {
					log.Fatalf("Error decoding item: %s", err)
				}
				if strings.Contains(loc.Data, "Africa") {
					count++
				}
			}
		default:
		}
	}

	runtime.GC()
	PrintMemUsage()

	fmt.Println("count =", count)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
