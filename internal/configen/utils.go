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
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

func resolve(env string, str string) (string, error) {
	t, err := template.New("str").Parse(str)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer

	err = t.Execute(&buff, Context{"Env": env})
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func mkdir(dir string) error {
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}

	index := filepath.Join(dir, "index.html")

	if _, err := os.Stat(index); os.IsNotExist(err) {
		if err := ioutil.WriteFile(index, []byte("<html></html>"), filePerm); err != nil {
			return err
		}
	}

	return nil
}

const (
	filePerm = 0600
	dirPerm  = 0755
)
