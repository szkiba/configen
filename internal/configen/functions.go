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
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
)

var functions = map[string]interface{}{
	"toYaml":        toYaml,
	"toToml":        toToml,
	"fromYaml":      fromYaml,
	"fromYamlArray": fromYamlArray,
	"fromToml":      fromToml,
	"equal":         reflect.DeepEqual,
	"assert":        assertion,
}

// ErrMissingValue returned by 'required' template function when required value is missing.
var ErrMissingValue = errors.New("missing value")

// ErrAssertionFailed returned by 'assert' template function when provided arg evaluated as flase value.
var ErrAssertionFailed = errors.New("assertion failed")

func toYaml(v interface{}) string {
	output, _ := yamlMarshal(v)

	return strings.TrimSuffix(string(output), "\n")
}

func toToml(v interface{}) string {
	output, _ := toml.Marshal(v)

	return strings.TrimSuffix(string(output), "\n")
}

func fromYaml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	yaml.Unmarshal([]byte(str), &m) // nolint

	return m
}

func fromYamlArray(str string) []interface{} {
	a := []interface{}{}

	yaml.Unmarshal([]byte(str), &a) // nolint

	return a
}

func fromToml(str string) map[string]interface{} {
	m := map[string]interface{}{}

	toml.Unmarshal([]byte(str), &m) // nolint

	return m
}

func assertion(msg string, v interface{}) (bool, error) {
	val := reflect.ValueOf(v)

	var valid bool

	switch val.Kind() { // nolint:exhaustive
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		valid = val.Len() != 0
	case reflect.Bool:
		valid = val.Bool()
	default:
		valid = val.IsValid() && !val.IsZero() && !val.IsNil()
	}

	if valid {
		return true, nil
	}

	return false, fmt.Errorf("%w: %s", ErrAssertionFailed, msg)
}

func newFuncMap(t *template.Template) template.FuncMap {
	funcs := sprig.TxtFuncMap()

	for k, v := range functions {
		funcs[k] = v
	}

	funcs["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		err := t.ExecuteTemplate(&buf, name, data)

		return buf.String(), err
	}

	funcs["required"] = func(msg string, v interface{}) (interface{}, error) {
		if v == nil {
			return "", fmt.Errorf("%w: %s", ErrMissingValue, msg)
		}

		if s, ok := v.(string); ok && s == "" {
			return "", fmt.Errorf("%w: %s", ErrMissingValue, msg)
		}

		return v, nil
	}

	return funcs
}

type files struct{}

func (f *files) Get(name string) (string, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (f *files) GetBytes(name string) ([]byte, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	return b, nil
}
