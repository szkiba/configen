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
	"fmt"
	"path/filepath"
	"strings"
)

const (
	templateLineNumberIdx = 2
)

func wrap(err error, path ...string) error {
	full := filepath.Join(path...)

	if !filepath.IsAbs(full) {
		if p, err := filepath.Abs(full); err == nil {
			full = p
		}
	}

	str := err.Error()

	if strings.HasPrefix(str, "template:") {
		if s := strings.Split(str, ":"); len(s) > templateLineNumberIdx {
			full += ":" + s[templateLineNumberIdx]
		}
	} else if strings.HasPrefix(str, "yaml:") {
		if s := strings.Split(str, ":"); len(s) > 1 && strings.HasPrefix(s[1], " line") {
			full += ":" + strings.TrimPrefix(s[1], " line ")
		}
	}

	return fmt.Errorf("%s: %w", full, err)
}
