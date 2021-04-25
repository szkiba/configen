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
	"github.com/antonmedv/expr"
	"github.com/itchyny/gojq"
	"github.com/jmespath/go-jmespath"
	"github.com/pelletier/go-toml"
	"github.com/qri-io/jsonpointer"
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
	"required":      required,
	"jq":            jq,
	"jp":            jp,
	"jptr":          jptr,
	"uritpl":        uritpl,
	"qsParse":       qsParse,
	"qsJoin":        qsJoin,
	"expr":          expression,
}

// ErrMissingValue returned by 'required' template function when required value is missing.
var ErrMissingValue = errors.New("missing value")

// ErrAssertionFailed returned by 'assert' template function when provided arg evaluated as flase value.
var ErrAssertionFailed = errors.New("assertion failed")

// ErrInvalidArgument returned when at least on argument from two is not a string.
var ErrInvalidArgument = errors.New("one of arguments should be string")

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

func required(msg string, v interface{}) (interface{}, error) {
	if v == nil {
		return "", fmt.Errorf("%w: %s", ErrMissingValue, msg)
	}

	if s, ok := v.(string); ok && s == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingValue, msg)
	}

	return v, nil
}

func newFuncMap() template.FuncMap {
	funcs := sprig.TxtFuncMap()

	for k, v := range functions {
		funcs[k] = v
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

func getStringAndInterface(a, b interface{}) (s string, v interface{}, err error) {
	if s1, ok := a.(string); ok {
		s = s1
		v = b
	} else if s2, ok := b.(string); ok {
		s = s2
		v = a
	} else {
		return "", nil, fmt.Errorf("%w: got %T, %T", ErrInvalidArgument, a, b)
	}

	if m, ok := v.(Context); ok {
		v = map[string]interface{}(m)
	}

	return s, v, nil
}

func jq(a, b interface{}) (interface{}, error) {
	str, v, err := getStringAndInterface(a, b)
	if err != nil {
		return nil, err
	}

	q, err := gojq.Parse(str)
	if err != nil {
		return nil, err
	}

	var val interface{}

	iter := q.Run(v)

	for {
		v, ok := iter.Next()
		if !ok {
			break
		}

		if e, ok := v.(error); ok {
			err = e

			break
		}

		val = v
	}

	if err != nil {
		return nil, err
	}

	return val, nil
}

func jp(a, b interface{}) (interface{}, error) {
	str, v, err := getStringAndInterface(a, b)
	if err != nil {
		return nil, err
	}

	res, err := jmespath.Search(str, v)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func jptr(a, b interface{}) (interface{}, error) {
	str, v, err := getStringAndInterface(a, b)
	if err != nil {
		return nil, err
	}

	ptr, err := jsonpointer.Parse(str)
	if err != nil {
		return nil, err
	}

	return ptr.Eval(v)
}

func expression(a, b interface{}) (interface{}, error) {
	str, v, err := getStringAndInterface(a, b)
	if err != nil {
		return nil, err
	}

	return expr.Eval(str, v)
}
