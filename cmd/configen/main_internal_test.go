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

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type args []string

func TestRun(t *testing.T) {
	t.Parallel()

	assert.Nil(t, os.Chdir("testdata"))

	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "YAML values", args: args{"-q", "--dump", "@dev", "format=json"}},
		{name: "TOML values", args: args{"-q", "--dump", "@test", "+values.toml", "format=json"}},
		{name: "JSON values", args: args{"-q", "--dump", "@demo", "+values.json", "format=json"}},
		{name: "same format", args: args{"-q", "--dump", "@sandbox"}},
		{name: "TOML format", args: args{"-q", "--dump", "@local", "format=toml"}},
		{name: "YAML format", args: args{"-q", "--dump", "@unstable", "format=yaml"}},

		{name: "help", args: args{"--help"}, want: 1},
		{name: "version", args: args{"--version"}, want: 0},
		{name: "error", args: args{"@unknown", "+missing"}, want: 1},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, run(append([]string{"configen"}, tt.args...)))
		})
	}
}
