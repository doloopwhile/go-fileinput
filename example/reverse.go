package main

import (
	"fmt"
	"os"

	"github.com/doloopwhile/go-fileinput"
)

func main() {
	sc := fileinput.Lines(os.Args[1:])
	lines := []string{}
	for sc.Scan() {
		l := sc.Text()
		lines = append(lines, l)
	}
	if sc.Err() != nil {
		os.Stderr.WriteString(sc.Err().Error() + "\n")
		os.Exit(1)
	}
	for i := len(lines) - 1; i >= 0; i-- {
		fmt.Printf("%02d: %s\n", i+1, lines[i])
	}
}
