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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/szkiba/configen/internal/configen"
)

func Test_newOptions(t *testing.T) {
	t.Parallel()

	if dir, _ := os.Getwd(); filepath.Base(dir) != "testdata" {
		assert.Nil(t, os.Chdir("testdata"))
	}

	tests := []struct {
		name    string
		args    args
		want    *options
		wantErr bool
	}{
		{
			name: "normal",
			want: &options{
				Options: configen.Options{ // nolint
					Templates: []string{"templates"},
					Output:    "dist",
					Values:    []string{"values.yaml"},
					Schemas:   []string{"schemas"},
					Package:   "package.json",
					Define:    make(map[string]string),
				},
				meta: meta{Env: []string{""}}, // nolint
			},
		},
		{
			name: "tags", args: args{"exe", "foo", "+values.json", "@test", "@dev"},
			want: &options{
				Options: configen.Options{ // nolint
					Templates: []string{"foo"},
					Output:    "dist/{{.Env}}",
					Values:    []string{"values.json"},
					Schemas:   []string{"schemas"},
					Package:   "package.json",
					Define:    make(map[string]string),
				},
				meta: meta{Env: []string{"test", "dev"}}, // nolint
			},
		},
		{name: "version", args: args{"--version"}, want: &options{Options: configen.Options{}, meta: meta{Version: true}}}, // nolint
		{name: "invalid dir", args: args{"--dir", "no such dir"}, wantErr: true},
		{name: "invalid flag", args: args{"--env", "--version"}, wantErr: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := newOptions(tt.args)

			if tt.wantErr {
				assert.Error(t, err)

				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
