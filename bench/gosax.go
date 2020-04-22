package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/eliben/gosax"
)

func main() {
	counter := 0
	inLocation := false

	PrintMemUsage()

	scb := gosax.SaxCallbacks{
		StartElement: func(name string, attrs []string) {
			if name == "location" {
				inLocation = true
			} else {
				inLocation = false
			}
		},

		EndElement: func(name string) {
			inLocation = false
		},

		Characters: func(contents string) {
			if inLocation && strings.Contains(contents, "Africa") {
				counter++
			}
		},
	}

	err := gosax.ParseFile(os.Args[1], scb)
	if err != nil {
		panic(err)
	}

	runtime.GC()
	PrintMemUsage()

	fmt.Println("counter =", counter)
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
