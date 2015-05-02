// An analogous article of fileinput module in Python
//
// Examples:
//
//    // reverse.go
//    // print file lines in reverse order.
//    package main
//
//    import (
//    	"fmt"
//    	"os"
//
//    	"github.com/doloopwhile/go-fileinput"
//    )
//
//    func main() {
//    	sc := fileinput.Lines(os.Args[1:])
//    	lines := []string{}
//    	for sc.Scan() {
//    		l := sc.Text()
//    		lines = append(lines, l)
//    	}
//    	if sc.Err() != nil {
//    		os.Stderr.WriteString(sc.Err().Error() + "\n")
//    		os.Exit(1)
//    	}
//    	for i := len(lines) - 1; i >= 0; i-- {
//    		fmt.Printf("%02d: %s\n", i+1, lines[i])
//    	}
//    }
package fileinput

import (
	"bufio"
	"io"
	"os"
)

// Lines returns a new Scanner to read lines of files in args.
// If args is empty, it return a Scanner which scans os.Stdin.
func Lines(args []string) *Scanner {
	if len(args) == 0 {
		args = []string{"-"}
	}
	return &Scanner{
		Args: args,
		Open: StdOpen,
	}
}

type (
	// Scanner provides a interface like bufio.Scanner
	// to reading data from multiple files.
	//
	// It is not expected that members of Scanner is modified after first call of .Scan()
	// If it was, it is undefined what happen.
	Scanner struct {
		Args      []string                                 // Names of files. It should be os.Args[1:] in typical use case.
		Open      func(name string) (io.ReadCloser, error) // Function to open files.
		SplitFunc bufio.SplitFunc                          // Argument of Split() of bufio.Split.

		sc    *bufio.Scanner
		rc    io.ReadCloser
		icurr int
		err   error
	}
)

// Scan advances internal scanner to the next token like Scan method of bufio.Scanner.
// It automatically open/close files specified in Args.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}
	if s.sc != nil {
		r := s.sc.Scan()
		s.err = s.sc.Err()
		if r {
			return true
		}
		s.rc.Close()
		s.sc = nil
		s.icurr++
	}

	for s.icurr < len(s.Args) {
		s.rc, s.err = s.Open(s.Args[s.icurr])
		if s.err != nil {
			return false
		}
		s.sc = bufio.NewScanner(s.rc)
		if s.SplitFunc != nil {
			s.sc.Split(s.SplitFunc)
		}
		r := s.sc.Scan()
		s.err = s.sc.Err()
		if r {
			return true
		}
		s.icurr++
	}
	return false
}

// Text returns the most recent token generated by a call to Scan as a newly allocated string holding its bytes.
func (s *Scanner) Text() string {
	if s.err != nil || s.sc == nil {
		return ""
	}
	return s.sc.Text()
}

// Err returns the first non-EOF error
// that was encountered by the Scanner or was returned by Open.
func (s *Scanner) Err() error {
	return s.err
}

// StdOpen open file with os.Open. However, it returns os.Stdin for "-".
func StdOpen(name string) (io.ReadCloser, error) {
	if name == "-" {
		return os.Stdin, nil
	}
	return os.Open(name)
}