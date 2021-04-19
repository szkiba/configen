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
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ErrValidationError returned if JSON schema validation failed.
var ErrValidationError = errors.New("validation error")

func (g *generator) validate(b []byte, format string) error {
	fn, ok := parsers[format]
	if !ok {
		return nil
	}

	v := Context{}

	if err := fn(b, &v); err != nil {
		return err
	}

	schema, ok := v.get(propSchema)
	if !ok || g.loose {
		return nil
	}

	loader := gojsonschema.NewGoLoader(v)

	result, err := gojsonschema.Validate(g.schemaLoader(schema), loader)
	if err != nil {
		return wrap(err, schema)
	}

	if !result.Valid() {
		for _, desc := range result.Errors() {
			fmt.Fprintf(os.Stderr, "%s\n", desc)
		}

		return ErrValidationError
	}

	return err
}

func (g *generator) schemaLoader(schema string) gojsonschema.JSONLoader {
	if l, ok := g.loaders[schema]; ok {
		return l
	}

	u, err := url.Parse(schema)
	if err == nil && u.Scheme == "" {
		u.Path = strings.TrimPrefix(path.Clean(u.Path), ".")
		if !strings.HasPrefix(u.Path, "/") {
			dir, _ := os.Getwd()
			u.Path = path.Join(filepath.ToSlash(dir), u.Path)
		}

		u.Scheme = "file"
		schema = u.String()
	}

	return gojsonschema.NewReferenceLoader(schema)
}
