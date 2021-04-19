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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerator_newContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     string
		values  []string
		want    Context
		wantErr bool
	}{
		{name: "normal", values: []string{"testdata/values/extra.yaml"}, env: "", want: extra},
		{name: "merged", values: []string{"testdata/values/extra.yaml", "testdata/values/values.yaml"}, env: "", want: merged}, // nolint:lll
		{name: "invalid name", values: []string{"{{env"}, env: "test", wantErr: true},                                          // nolint:lll
		{name: "invalid content", values: []string{"testdata/values/badextra.json"}, env: "", wantErr: true},                   // nolint:lll
		{name: "missing file", values: []string{"no such file"}, env: "", wantErr: true},                                       // nolint:lll
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			g := new(generator)
			o := new(Options)
			o.Values = tt.values
			o.Define = map[string]string{}

			c, err := g.newContext(tt.env, o)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Nil(t, err)
			assert.EqualValues(t, tt.want, c["Values"])
			assert.IsType(t, &files{}, c["Files"])
			assert.Equal(t, tt.env, c["Env"])
		})
	}
}
