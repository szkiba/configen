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
	"encoding/json"
	"errors"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"
	"muzzammil.xyz/jsonc"
)

// Context defines generic JSON/YAML/TOML values type.
type Context map[string]interface{}

// ErrUnknownFormat returned when file format unsupported or unrecognizable from file extension.
var ErrUnknownFormat = errors.New("unknown format")

type parseFunc func([]byte, interface{}) error

type formatFunc func(interface{}) ([]byte, error)

var (
	parsers = map[string]parseFunc{
		"yaml":  yaml.Unmarshal,
		"yml":   yaml.Unmarshal,
		"json":  jsonUnmarshal,
		"jsonc": jsonUnmarshal,
		"toml":  toml.Unmarshal,
	}

	formatters = map[string]formatFunc{
		"yaml":  yamlMarshal,
		"yml":   yamlMarshal,
		"json":  jsonMarshal,
		"jsonc": jsonMarshal,
		"toml":  toml.Marshal,
	}
)

const (
	propSchema = "$schema"
	propFormat = "$format"
	yamlIndent = 2
)

func jsonMarshal(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return jsonc.Unmarshal(data, v)
}

func yamlMarshal(data interface{}) ([]byte, error) {
	var buff bytes.Buffer

	enc := yaml.NewEncoder(&buff)
	enc.SetIndent(yamlIndent)

	if err := enc.Encode(data); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (c *Context) unmarshal(data []byte, format string) error {
	fn, ok := parsers[format]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownFormat, format)
	}

	return fn(data, c)
}

func (c Context) marshal(format string) ([]byte, error) {
	fn, ok := formatters[format]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrUnknownFormat, format)
	}

	return fn(c)
}

func (c Context) merge(other Context) error {
	return mergo.Merge(&c, other)
}

func (c Context) get(key string) (string, bool) {
	val, ok := c[key]
	if !ok {
		return "", false
	}

	s, ok := val.(string)
	if !ok {
		return "", false
	}

	return s, true
}
