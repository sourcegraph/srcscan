package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/sourcegraph/srcscan"
	"os"
	"reflect"
)

var verbose = flag.Bool("v", false, "show verbose output")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: srcscan [OPTS] DIR..\n\n")
		fmt.Fprintf(os.Stderr, "where OPTS is any of:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	flag.Parse()

	var dirs []string
	if flag.NArg() == 0 {
		dirs = []string{"."}
	} else {
		dirs = flag.Args()
	}

	for i, dir := range dirs {
		units, err := srcscan.Scan(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
		for j, unit := range units {
			fmt.Printf("%-15s %s\n", typeName(unit), unit.Path())
			if *verbose {
				out, err := json.MarshalIndent(unit, "    ", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "error serializing to JSON: %s\n", err)
					os.Exit(1)
				}
				fmt.Printf("    %s\n", out)
				if i != len(dirs)-1 || j != len(units)-1 {
					fmt.Printf("\n")
				}
			}
		}
	}
}

func typeName(v interface{}) string {
	return reflect.TypeOf(v).Elem().Name()
}
