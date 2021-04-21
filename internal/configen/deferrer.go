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
	"os"
	"text/template"
)

type deferrer struct {
	quiet    bool
	context  Context
	template *template.Template
	deferred []string
}

func newDeferrer(quiet bool, t *template.Template, ctx Context) *deferrer {
	d := &deferrer{
		quiet:    quiet,
		context:  ctx,
		template: t,
		deferred: []string{},
	}

	funcs := template.FuncMap{}

	funcs["defer"] = func(name string) string {
		d.deferred = append([]string{name}, d.deferred...)

		return name
	}

	t.Funcs(funcs)

	return d
}

func (d *deferrer) render(document interface{}) error {
	if len(d.deferred) == 0 {
		return nil
	}

	d.context["Document"] = document

	for _, name := range d.deferred {
		var buff bytes.Buffer

		if err := d.template.ExecuteTemplate(&buff, name, d.context); err != nil {
			return err
		}

		if !d.quiet {
			if _, err := os.Stdout.Write(buff.Bytes()); err != nil {
				return err
			}
		}
	}

	return nil
}
