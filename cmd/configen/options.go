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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/szkiba/configen/internal/configen"
)

const (
	app  = "configen"
	desc = `Template based configuration generator.

You can specify multiple environments, input directories and values files.
Frequently used options has alternative positional argument syntax.`
)

type meta struct {
	Env     []string `short:"e" long:"env" value-name:"environment" description:"Staging environment name [arg: @environment]"` //nolint:lll
	Dir     string   `long:"dir" value-name:"directory" description:"Set working directory"`
	Version bool     `short:"V" long:"version" description:"Show version information"`
}

type options struct {
	configen.Options
	meta
}

func newOptions(args []string) (*options, error) {
	opts := new(options)

	_, err := flags.NewParser(&opts.meta, flags.IgnoreUnknown|flags.PrintErrors).ParseArgs(args)
	if err != nil {
		return nil, err
	}

	if opts.meta.Version {
		return opts, nil
	}

	parser := flags.NewNamedParser(app, flags.Default)
	parser.Usage = "[options] [args]"
	parser.Command.Group.LongDescription = desc

	if _, err = parser.AddGroup("Options", "", opts); err != nil {
		return nil, err
	}

	positional, err := parser.ParseArgs(args)
	if err != nil {
		return nil, err
	}

	if err := opts.setWorkingDir(); err != nil {
		return nil, err
	}

	opts.applyTags(positional)
	opts.applyDefaults()

	return opts, err
}

func (o *options) setWorkingDir() error {
	if len(o.Dir) > 0 {
		if err := os.Chdir(o.Dir); err != nil {
			fmt.Fprintln(os.Stderr, err)

			return err
		}
	}

	return nil
}

func (o *options) defaultTemplate() []string {
	return []string{"templates"}
}

func (o *options) defaultSchemas() []string {
	if info, err := os.Stat("schemas"); err == nil && info.IsDir() {
		return []string{"schemas"}
	}

	return []string{}
}

func (o *options) defaultPackage() string {
	if info, err := os.Stat("package.json"); err == nil && !info.IsDir() {
		return "package.json"
	}

	return ""
}

func (o *options) defaultOutput() string {
	if len(o.Env) == 0 {
		return "dist"
	}

	return filepath.Join("dist", "{{.Env}}")
}

func (o *options) defaultValues() []string {
	if val := findValues(); len(val) > 0 {
		return []string{val}
	}

	return []string{}
}

func (o *options) applyDefaults() {
	if len(o.Templates) == 0 {
		o.Templates = o.defaultTemplate()
	}

	if len(o.Schemas) == 0 {
		o.Schemas = o.defaultSchemas()
	}

	if len(o.Package) == 0 {
		o.Package = o.defaultPackage()
	}

	if len(o.Output) == 0 {
		o.Output = o.defaultOutput()
	}

	if len(o.Values) == 0 {
		o.Values = o.defaultValues()
	}

	if len(o.Env) == 0 {
		o.Env = []string{""}
	}
}

func (o *options) applyTags(args []string) {
	if len(args) > 1 {
		for _, t := range args[1:] {
			o.applyTag(t)
		}
	}
}

func (o *options) applyTag(t string) {
	if strings.HasPrefix(t, "@") {
		o.Env = append(o.Env, t[1:])

		return
	}

	if strings.HasPrefix(t, "+") {
		o.Values = append(o.Values, t[1:])

		return
	}

	if strings.ContainsRune(t, '=') {
		f := strings.SplitN(t, "=", 2)
		o.Define[f[0]] = f[1]

		return
	}

	o.Templates = append(o.Templates, t)
}

var defaultValues = []string{"values.yaml", "values.yml", "values.toml", "values.json"}

func findValues() string {
	for _, name := range defaultValues {
		if info, err := os.Stat(name); err == nil && !info.IsDir() {
			return name
		}
	}

	return ""
}
