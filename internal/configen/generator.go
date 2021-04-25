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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobwas/glob"
	"github.com/jpillora/longestcommon"
	"github.com/xeipuuv/gojsonschema"
)

// Generate is the main entry point, called after parsing command line.
func Generate(opts *Options, envs ...string) error {
	dirs := make([]string, len(envs)+1)

	dirs[0] = opts.Output

	for i, env := range envs {
		g, err := newGenerator(env, opts)
		if err != nil {
			return err
		}

		if err := g.generateDir(opts.Templates...); err != nil {
			return err
		}

		dirs[1+i] = g.output
	}

	if !opts.Dry && len(opts.Package) > 0 {
		return preparePackage(longestcommon.Prefix(dirs), opts.Package)
	}

	return nil
}

func preparePackage(dir string, file string) error {
	mkdir(dir) // nolint

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	out := filepath.Join(dir, filepath.Base(file))

	return ioutil.WriteFile(out, b, filePerm)
}

type generator struct {
	output  string
	loaders schemaLoaders
	dump    bool
	loose   bool
	dry     bool
	quiet   bool
	ctx     Context
	root    *template.Template
}

func newGenerator(env string, o *Options) (g *generator, err error) {
	g = new(generator)

	g.dump = o.Dump
	g.loose = o.Loose
	g.dry = o.Dry
	g.quiet = o.Quiet

	if g.root, err = g.newRootTemplate(env, o); err != nil {
		return nil, err
	}

	if g.loaders, err = g.newSchemaLoader(env, o); err != nil {
		return nil, err
	}

	if g.ctx, err = g.newContext(env, o); err != nil {
		return nil, err
	}

	if g.output, err = resolve(env, o.Output); err != nil {
		return nil, err
	}

	return g, nil
}

func (g *generator) generateDir(dirs ...string) error {
	for _, dir := range dirs {
		dir := dir
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() || partialGlobe.Match(path) {
					return nil
				}

				rel, err := filepath.Rel(dir, path)
				if err != nil {
					return err
				}

				return g.generateFile(dir, rel)
			})
		// nolint
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *generator) templateFuncMap(t *template.Template) template.FuncMap {
	funcs := g.newFuncMap()

	funcs["include"] = func(name string, data interface{}) (string, error) {
		var buf strings.Builder
		err := t.ExecuteTemplate(&buf, name, data)

		return buf.String(), err
	}

	funcs["tpl"] = func(tpl string, ctx Context) (string, error) {
		tmpl, err := t.New("").Parse(tpl)
		if err != nil {
			return "", err
		}

		sub := ctx

		if err := sub.merge(g.ctx); err != nil {
			return "", err
		}

		var buff bytes.Buffer

		if err := tmpl.Execute(&buff, sub); err != nil {
			return "", err
		}

		return buff.String(), nil
	}

	return funcs
}

func (g *generator) newFuncMap() template.FuncMap {
	funcs := newFuncMap()

	funcs["validate"] = func(schema string, v map[string]interface{}) bool {
		err := g.validate(schema, v)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)

			return false
		}

		return true
	}

	if g.quiet {
		noop := func(a ...interface{}) (int, error) {
			return 0, nil
		}
		funcs["out"] = noop
		funcs["outln"] = noop
		funcs["outf"] = noop
	} else {
		funcs["out"] = func(a ...interface{}) (int, error) {
			return fmt.Fprint(os.Stdout, a...) // nolint
		}
		funcs["outln"] = func(a ...interface{}) (int, error) {
			return fmt.Fprintln(os.Stdout, a...) // nolint
		}
		funcs["outf"] = func(format string, a ...interface{}) (int, error) {
			return fmt.Fprintf(os.Stdout, format, a...) // nolint
		}
	}

	return funcs
}

func (g *generator) executeTemplate(basedir string, path string) ([]byte, *deferrer, error) {
	src := filepath.Join(basedir, path)

	t, err := g.root.Clone()
	if err != nil {
		return nil, nil, wrap(err, src)
	}

	t.Funcs(g.templateFuncMap(t))

	ctx := g.ctx
	def := newDeferrer(g.quiet, t, g.ctx)

	t, err = t.ParseFiles(src)
	if err != nil {
		return nil, nil, wrap(err, src)
	}

	var buff bytes.Buffer

	err = t.ExecuteTemplate(&buff, filepath.Base(src), ctx)
	if err != nil {
		return nil, nil, wrap(err, src)
	}

	txt := buff.Bytes()

	if g.dump {
		dump := filepath.Join(g.output, path) + dumpSuffix

		dir := filepath.Dir(dump)

		if err := mkdir(dir); err != nil {
			return nil, nil, wrap(err, dir)
		}

		if err := ioutil.WriteFile(dump, txt, filePerm); err != nil {
			return nil, nil, wrap(err, dump)
		}
	}

	return txt, def, nil
}

func (g *generator) generateFile(basedir string, path string) error {
	txt, console, err := g.executeTemplate(basedir, path)
	if err != nil {
		return err
	}

	out := filepath.Join(g.output, path)

	errfile := filepath.Join(basedir, path)

	if g.dump {
		errfile = out + dumpSuffix
	}

	if !g.dry {
		dir := filepath.Dir(out)

		if err := mkdir(dir); err != nil {
			return wrap(err, dir)
		}
	}

	format := strings.TrimPrefix(filepath.Ext(out), ".")

	txt, format, err = transform(txt, format)
	if err != nil {
		return wrap(err, errfile)
	}

	out = outname(out, format)

	parsed, err := g.validateRaw(txt, format)
	if err != nil {
		return wrap(err, errfile)
	}

	if err := console.render(parsed); err != nil {
		return err
	}

	if g.dry {
		return nil
	}

	err = ioutil.WriteFile(out, txt, filePerm)
	if err != nil {
		return err
	}

	return nil
}

func outname(out, format string) string {
	if len(format) > 0 {
		return strings.TrimSuffix(out, filepath.Ext(out)) + "." + format
	}

	return out
}

func (g *generator) newRootTemplate(env string, o *Options) (*template.Template, error) {
	t := template.New(partialPrefix)

	t = t.Funcs(g.templateFuncMap(t))

	for _, dir := range o.Templates {
		dir, err := resolve(env, dir)
		if err != nil {
			return nil, wrap(err, dir)
		}

		err = filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() || !partialGlobe.Match(path) {
					return nil
				}

				t, err = t.ParseFiles(path)
				if err != nil {
					return wrap(err, path)
				}

				return nil
			})
		// nolint
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}

func (g *generator) newContext(env string, o *Options) (Context, error) {
	values := Context{}

	for name, value := range o.Define {
		values[name] = value
	}

	for _, file := range o.Values {
		f, err := resolve(env, file)
		if err != nil {
			return nil, wrap(err, file)
		}

		b, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, wrap(err, f)
		}

		format := strings.TrimPrefix(filepath.Ext(f), ".")

		ctx := Context{}
		if err := ctx.unmarshal(b, format); err != nil {
			return nil, wrap(err, f)
		}

		if _, err = g.validateRaw(b, format); err != nil {
			return nil, wrap(err, f)
		}

		if err := values.merge(ctx); err != nil {
			return nil, err
		}
	}

	return Context{"Values": values, "Files": &files{}, "Env": env}, nil
}

func (g *generator) newSchemaLoader(env string, o *Options) (schemaLoaders, error) {
	loaders := schemaLoaders{}

	for _, dir := range o.Schemas {
		dir, err := resolve(env, dir)
		if err != nil {
			return nil, wrap(err, dir)
		}

		err = filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() || filepath.Ext(path) != ".json" {
					return nil
				}

				b, err := ioutil.ReadFile(path)
				if err != nil {
					return wrap(err, path)
				}

				format := strings.TrimPrefix(filepath.Ext(path), ".")

				ctx := Context{}
				if err := ctx.unmarshal(b, format); err != nil {
					return wrap(err, path)
				}

				id, ok := ctx.get("$id")
				if !ok {
					return nil
				}

				loaders[id] = gojsonschema.NewGoLoader(ctx)

				return nil
			})
		// nolint
		if err != nil {
			return nil, err
		}
	}

	return loaders, nil
}

const (
	partialPrefix = "_"
	dumpSuffix    = "~"
)

var partialGlobe = glob.MustCompile(filepath.Join("**", "_*"))
