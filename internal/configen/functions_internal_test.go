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
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func Test_newFuncMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "sprig", in: `{{trim " Hello "}}`, want: "Hello"},
		{name: "toYaml", in: `{{toYaml .values}}`, want: strings.TrimSpace(string(valuesYAML))},
		{name: "toToml", in: `{{toToml .values}}`, want: strings.TrimSpace(string(valuesTOML))},
		{name: "fromYaml", in: `{{toYaml .values|fromYaml|toToml}}`, want: strings.TrimSpace(string(valuesTOML))},
		{name: "fromToml", in: `{{toToml .values|fromToml|toYaml}}`, want: strings.TrimSpace(string(valuesYAML))},
		{name: "fromYamlArray", in: `{{list "foo" "bar"|toYaml|fromYamlArray|join "+"}}`, want: "foo+bar"},

		{name: "jp", in: `{{jp "values.name" .}}`, want: "foo"},
		{name: "jp", in: `{{jp . "values.name"}}`, want: "foo"},

		{name: "jptr", in: `{{jptr . "/values/name"}}`, want: "foo"},

		{name: "expr", in: `{{expr . "values.name"}}`, want: "foo"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpl := template.New("test")
			funcs := newFuncMap()

			tmpl, err := tmpl.Funcs(funcs).Parse(tt.in)

			assert.Nil(t, err)

			var buff bytes.Buffer

			err = tmpl.Execute(&buff, Context{"values": values})

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.want, buff.String())
		})
	}
}

func Test_files(t *testing.T) {
	t.Parallel()

	f := &files{}

	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "normal", in: "testdata/values/values.toml", want: string(valuesTOML)},
		{name: "missing", in: "no such file", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str, err1 := f.Get(tt.in)
			b, err2 := f.GetBytes(tt.in)

			if tt.wantErr {
				assert.Error(t, err1)
				assert.Error(t, err2)

				return
			}

			assert.Nil(t, err1)
			assert.Nil(t, err2)
			assert.Equal(t, tt.want, str)
			assert.Equal(t, tt.want, string(b))
		})
	}
}
