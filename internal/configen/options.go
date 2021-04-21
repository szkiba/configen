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

// Options holds command line flags.
type Options struct {
	Templates []string          `short:"t" long:"template" value-name:"directory" description:"Input directory [arg: directory] (default: templates)"` //nolint:lll
	Output    string            `short:"o" long:"output" value-name:"directory" description:"Output directory (default: dist)"`                        //nolint:lll
	Schemas   []string          `short:"s" long:"schema" value-name:"directory" description:"Schema directory (default: schemas)"`                     //nolint:lll
	Values    []string          `short:"f" long:"values" value-name:"file" description:"Data values file [arg: +file] (default: values.yaml)"`         //nolint:lll
	Define    map[string]string `long:"set" value-name:"name:value" description:"Set value [arg: name=value]"`                                         //nolint:lll
	Loose     bool              `long:"loose" description:"Disable schema validation"`
	Dry       bool              `long:"dry-run" description:"Skip writing output files"`
	Dump      bool              `long:"dump" description:"Dump intermediate files"`
	Quiet     bool              `short:"q" long:"quiet" description:"Suppress console output"`
	Package   string            `short:"p" long:"package" value-name:"file" description:"Package descriptor template (default: package.json)"` //nolint:lll
}
