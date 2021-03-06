package main

// this is a manual copy of myitcv.io/_tmpls/cliflag to bootstrap pkgconcat

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	usage string
)

func setupAndParseFlags(msg string) {
	flag.Usage = func() {
		res := new(strings.Builder)
		fmt.Fprint(res, msg)

		res.WriteString("Flags:\n")
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
	var indent string // the indent (if there is one) when we carry

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
		sincelastTab := 0
		seenTab := false

		nc := len(indent)

		if nc >= width {
			fatalf("cannot foldOnSpaces where indent is greater than width")
		}

		var postSpace string
		var space string

	Space:
		for i := 0; i < len(carry); {
			prevPos := i
			nbs, rw := utf8.DecodeRuneInString(carry[i:])
			i += rw

			nc++

			if nbs == '\t' {
				seenTab = true
			}

			if !seenTab {
				sincelastTab++
			}

			if isSplitter(nbs) {
				if postSpace != "" {
					res.WriteString(space)
					res.WriteString(postSpace)
					space = string(nbs)
				} else {
					space += string(nbs)
				}

				if nc >= width {
					res.WriteString("\n")
					carry = strings.TrimLeftFunc(carry[prevPos:], unicode.IsSpace)

					if seenTab {
						indent += strings.Repeat(" ", sincelastTab) + "\t"
					}

					continue Line
				}

				lastSpace = nc
				postSpace = ""

			} else {
				if lastSpace == -1 {
					res.WriteRune(nbs)
					continue Space
				}

				if nc == width {
					if seenTab {
						indent += strings.Repeat(" ", sincelastTab) + "\t"
					}

					res.WriteString("\n")
					carry = strings.TrimLeftFunc(postSpace+carry[prevPos:], unicode.IsSpace)
					continue Line
				}

				postSpace += string(nbs)
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
func infof(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}
func info(args ...interface{}) {
	fmt.Fprint(os.Stderr, args...)
}
func infoln(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
