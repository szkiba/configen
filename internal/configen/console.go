// MIT License
//
// Copyright (c) 2021 Iván Szkiba
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package configen

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

type console struct {
	quiet    bool
	context  Context
	template *template.Template
	deferred []string
}

func newConsole(quiet bool, t *template.Template, ctx Context) *console {
	c := &console{
		quiet:    quiet,
		context:  ctx,
		template: t,
		deferred: []string{},
	}

	t.Funcs(c.funcMap())

	return c
}

func (c *console) funcMap() template.FuncMap {
	funcs := template.FuncMap{}

	if c.quiet {
		noop := func(a ...interface{}) (int, error) {
			return 0, nil
		}
		funcs["out"] = noop
		funcs["outln"] = noop
		funcs["outf"] = noop
	} else {
		funcs["out"] = func(a ...interface{}) (int, error) {
			return fmt.Fprint(os.Stdout, a...) // nolint
		}
		funcs["outln"] = func(a ...interface{}) (int, error) {
			return fmt.Fprintln(os.Stdout, a...) // nolint
		}
		funcs["outf"] = func(format string, a ...interface{}) (int, error) {
			return fmt.Fprintf(os.Stdout, format, a...) // nolint
		}
	}

	funcs["defer"] = func(name string) string {
		c.deferred = append([]string{name}, c.deferred...)

		return name
	}

	return funcs
}

func (c *console) render(document interface{}) error {
	if len(c.deferred) == 0 {
		return nil
	}

	c.context["Document"] = document

	for _, name := range c.deferred {
		var buff bytes.Buffer

		if err := c.template.ExecuteTemplate(&buff, name, c.context); err != nil {
			return err
		}

		if !c.quiet {
			if _, err := os.Stdout.Write(buff.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}
