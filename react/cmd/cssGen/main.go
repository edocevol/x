// cssGen is a temporary code generator for the myitcv.io/react.CSS type
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"
	"unicode/utf8"

	"myitcv.io/gogenerate"
)

// using
// https://github.com/Microsoft/TypeScript/blob/8b9fa4ce7420fdf2f540300dc80fa91f5b89ea93/lib/lib.dom.d.ts#L1692
// as a reference
//
var attrs = map[string]*typ{
	"Float":      &typ{},
	"FontSize":   &typ{HTML: "font-size"},
	"FontStyle":  &typ{HTML: "font-style"},
	"FontWeight": &typ{HTML: "font-weight"},
	"Height":     &typ{},
	"MaxHeight":  &typ{HTML: "max-height"},
	"MarginTop":  &typ{HTML: "margin-top"},
	"OverflowY":  &typ{HTML: "overflow-y"},
	"MinHeight":  &typ{HTML: "min-height"},
	"Overflow":   &typ{},
	"Resize":     &typ{},
	"Width":      &typ{},
	"Position":   &typ{},
	"Top":        &typ{},
	"Left":       &typ{},
	"ZIndex":     &typ{HTML: "z-index"},
}

const (
	cssGenCmd = "cssGen"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(cssGenCmd + ": ")

	flag.Parse()

	for n, a := range attrs {
		a.Name = n
		if a.React == "" {
			a.React = lowerInitial(n)
		}
		if a.HTML == "" {
			a.HTML = strings.ToLower(n)
		}
		if a.Type == "" {
			a.Type = "string"
		}
	}

	write := func(tmpl string, fn string) {
		buf := bytes.NewBuffer(nil)

		t, err := template.New("t").Parse(tmpl)
		if err != nil {
			fatalf("could not parse template: %v", err)
		}

		err = t.Execute(buf, attrs)
		if err != nil {
			fatalf("could not execute template: %v", err)
		}

		toWrite := buf.Bytes()
		out, err := format.Source(toWrite)
		if err == nil {
			toWrite = out
		}

		if err := ioutil.WriteFile(fn, toWrite, 0644); err != nil {
			fatalf("could not write %v: %v", fn, err)
		}
	}

	write(tmpl, gogenerate.NameFile("react", cssGenCmd))
	write(jsxTmpl, filepath.Join("jsx", gogenerate.NameFile("jsx", cssGenCmd)))
}

func lowerInitial(s string) string {
	if s == "" {
		return ""
	}

	r, w := utf8.DecodeRuneInString(s)
	return strings.ToLower(string(r)) + s[w:]
}

type typ struct {
	Name string

	// React is the React property name if not equivalent to the lower-initial
	// camel-case version of .Name
	React string

	// HTML is the HTML property name if not equivalent to the lowercase version
	// of .Name
	HTML string

	// Type is the type. Default is "string"
	Type string
}

var tmpl = `
 // Code generated by cssGen. DO NOT EDIT.

package react

import "github.com/gopherjs/gopherjs/js"

// CSS defines CSS attributes for HTML components. Largely based on
// https://developer.mozilla.org/en-US/docs/Web/CSS/Reference
//
type CSS struct {
	o *js.Object

	{{range . }}
	{{.Name}} {{.Type}}
	{{- end}}
}

// TODO: until we have a resolution on
// https://github.com/gopherjs/gopherjs/issues/236 we define hack() below

func (c *CSS) hack() *CSS {
	if c == nil {
		return nil
	}

	o := object.New()

	{{range . }}
	o.Set("{{.React}}", c.{{.Name}})
	{{- end}}

	return &CSS{o: o}
}
`

var jsxTmpl = `
package jsx

import (
	"fmt"
	"strings"

	"myitcv.io/react"
)

func parseCSS(s string) *react.CSS {
	res := new(react.CSS)

	parts := strings.Split(s, ";")

	for _, p := range parts {
		kv := strings.Split(p, ":")
		if len(kv) != 2 {
			panic(fmt.Errorf("invalid key-val %q in %q", p, s))
		}

		k, v := kv[0], kv[1]

		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, "\"")

		switch k {
		{{range .}}
		case "{{.HTML}}":
			res.{{.Name}} = v
		{{end}}
		default:
			panic(fmt.Errorf("unknown CSS key %q in %q", k, s))
		}
	}

	return res
}
`

func fatalf(format string, args ...interface{}) {
	panic(fmt.Errorf(format, args...))
}
