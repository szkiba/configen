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
	_ "embed" // nolint
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	values = map[string]interface{}{"name": "foo", "version": "1.0.0"}
	merged = map[string]interface{}{"name": "foo", "version": "1.0.0", "extra": "bar"}
	extra  = map[string]interface{}{"extra": "bar"}

	valuesYAML = []byte(`name: foo
version: 1.0.0
`)

	valuesJSON = []byte(`{
  "name": "foo",
  "version": "1.0.0"
}`)

	valuesTOML = []byte(`name = "foo"
version = "1.0.0"
`)
)

func TestContext_unmarshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		format  string
		wantErr bool
	}{
		{name: "yaml", data: valuesYAML, format: "yaml"},
		{name: "yml", data: valuesYAML, format: "yml"},
		{name: "json", data: valuesJSON, format: "json"},
		{name: "toml", data: valuesTOML, format: "toml"},

		{name: "unknown format", data: valuesYAML, format: "unknown", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := Context{}

			err := c.unmarshal(tt.data, tt.format)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Nil(t, err)
			assert.EqualValues(t, values, c)
		})
	}
}

func TestContext_marshal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		format  string
		wantErr bool
	}{
		{name: "yaml", data: valuesYAML, format: "yaml"},
		{name: "yml", data: valuesYAML, format: "yml"},
		{name: "json", data: valuesJSON, format: "json"},
		{name: "toml", data: valuesTOML, format: "toml"},

		{name: "unknown format", format: "unknown", wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := Context(values)

			bytes, err := c.marshal(tt.format)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Nil(t, err)
			assert.EqualValues(t, tt.data, bytes)
		})
	}
}

func TestContext_get(t *testing.T) {
	t.Parallel()

	c := Context(values)

	val, ok := c.get("name")

	assert.True(t, ok)
	assert.Equal(t, "foo", val)

	val, ok = c.get("missing value")

	assert.False(t, ok)
	assert.Equal(t, "", val)

	c = Context{"name": true}

	val, ok = c.get("name")

	assert.False(t, ok)
	assert.Equal(t, "", val)
}
