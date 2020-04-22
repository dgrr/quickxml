// +build ignore
package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	xml "github.com/dgrr/quickxml"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var (
		readNext = false
		count    = 0
	)

	PrintMemUsage()

	r := xml.NewReader(file)
	for r.Next() {
		switch e := r.Element().(type) {
		case *xml.StartElement:
			readNext = e.NameUnsafe() == "location"
		case *xml.TextElement:
			if readNext && strings.Contains(string(*e), "Africa") {
				count++
				readNext = false
			}
		}
	}

	runtime.GC()
	PrintMemUsage()

	fmt.Println("counter =", count)
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
