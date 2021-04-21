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
	"net/url"
	"reflect"

	"github.com/yosida95/uritemplate/v3"
)

func qsParse(str string) (map[string][]string, error) {
	return url.ParseQuery(str)
}

func qsJoin(v map[string][]string) string {
	return url.Values(v).Encode()
}

func uritpl(tmpl string, params map[string]interface{}) (string, error) {
	t, err := uritemplate.New(tmpl)
	if err != nil {
		return "", err
	}

	vars := uritemplate.Values{}

	for key, val := range params {
		vars.Set(key, uritplValue(val))
	}

	return t.Expand(vars)
}

func uritplValue(val interface{}) uritemplate.Value { // nolint: cyclop
	if s, ok := val.(string); ok && len(s) > 0 {
		return uritemplate.String(s)
	}

	if l, ok := val.([]string); ok && len(l) > 0 {
		return uritemplate.List(l...)
	}

	if q, ok := val.(map[string][]string); ok && len(q) > 0 {
		return uritplParams(q)
	}

	if q, ok := val.(url.Values); ok && len(q) > 0 {
		return uritplParams(q)
	}

	if m, ok := val.(map[string]interface{}); ok && len(m) > 0 {
		return uritplMap(m)
	}

	if m, ok := val.(map[string]string); ok && len(m) > 0 {
		return uritplStringMap(m)
	}

	ref := reflect.ValueOf(val)

	switch ref.Kind() { // nolint: exhaustive
	case reflect.Map, reflect.Array, reflect.Slice, reflect.String:
		if ref.Len() == 0 {
			return uritemplate.KV()
		}
	}

	return uritemplate.String(fmt.Sprint(val))
}

func uritplMap(m map[string]interface{}) uritemplate.Value {
	kv := make([]string, 0, len(m)<<1)

	for k, v := range m {
		kv = append(kv, k, v.(string))
	}

	return uritemplate.KV(kv...)
}

func uritplStringMap(m map[string]string) uritemplate.Value {
	kv := make([]string, 0, len(m)<<1)

	for k, v := range m {
		kv = append(kv, k, v)
	}

	return uritemplate.KV(kv...)
}

func uritplParams(q map[string][]string) uritemplate.Value {
	kv := make([]string, 0, len(q)<<1)

	for k, a := range q {
		for _, v := range a {
			kv = append(kv, k, v)
		}
	}

	return uritemplate.KV(kv...)
}
