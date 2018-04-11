// Code generated by pkgconcat. DO NOT EDIT.

package main
import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
	"golang.org/x/text/unicode/norm"
)
var (
	usage string
)
func setupAndParseFlags(msg string) {
	flag.Usage = func() {
		res := new(strings.Builder)
		fmt.Fprint(res, msg)

		flag.CommandLine.SetOutput(res)
		flag.PrintDefaults()
		res.WriteString("\n")
		res.WriteString("\n")

		fmt.Fprint(os.Stderr, foldOnSpaces(res.String(), 80))

		os.Exit(0)
	}
	flag.Parse()

	flag.CommandLine.SetOutput(os.Stderr)
}
func foldOnSpaces(input string, width int) string {
	var carry string
	var indent string	// the indent (if there is one) when we carry

	sc := bufio.NewScanner(strings.NewReader(input))

	res := new(strings.Builder)
	first := true

Line:
	for {
		carried := carry != ""
		if !carried {
			if !sc.Scan() {
				break
			}

			if first {
				first = false
			} else {
				res.WriteString("\n")
			}

			carry = sc.Text()

			iBuilder := new(strings.Builder)

			for _, r := range carry {
				if !unicode.IsSpace(r) {
					break
				}
				iBuilder.WriteRune(r)
			}

			indent = iBuilder.String()

			carry = strings.TrimSpace(carry)
		}

		if len(carry) == 0 {
			continue
		}

		res.WriteString(indent)

		if len(indent)+len(carry) < width {
			res.WriteString(carry)
			carry = ""
			continue
		}

		lastSpace := -1

		var ia norm.Iter
		ia.InitString(norm.NFD, carry)
		nc := len(indent)

		if nc >= width {
			fatalf("cannot foldOnSpaces where indent is greater than width")
		}

		var postSpace string

	Space:
		for !ia.Done() {
			prevPos := ia.Pos()
			nbs := ia.Next()
			r, rw := utf8.DecodeRune(nbs)
			if rw != len(nbs) {
				fatalf("didn't expect a multi-rune normalisation response: %v", string(nbs))
			}

			nc++

			spaceCount := 0

			if isSplitter(r) {
				spaceCount++
			}

			switch spaceCount {
			case 0:

				if lastSpace == -1 {
					res.WriteRune(r)
					continue Space
				}

				if nc == width {

					res.WriteString("\n")
					carry = strings.TrimLeftFunc(postSpace+carry[prevPos:], unicode.IsSpace)
					continue Line
				}

				postSpace += string(r)
				continue Space
			case 1:

				res.WriteString(postSpace)

				switch {
				case nc == width:
					res.WriteRune(r)
					fallthrough
				case nc > width:
					res.WriteString("\n")
					carry = strings.TrimLeftFunc(carry[ia.Pos():], unicode.IsSpace)

					continue Line
				}

				res.WriteRune(r)

				lastSpace = nc
				postSpace = ""
				continue Space
			default:
				fatalf("is this even possible?")
			}
		}

		carry = ""
	}

	if err := sc.Err(); err != nil {
		fatalf("failed to scan in foldOnSpaces: %v", err)
	}

	return res.String()
}
func isSplitter(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}

	switch r {
	case '/':
		return true
	}

	return false
}
func fatalf(format string, args ...interface{}) {
	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
