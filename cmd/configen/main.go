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
	"runtime"

	"github.com/szkiba/configen/internal/configen"
)

var version = "dev"

func run(args []string) int {
	opts, err := newOptions(args)
	if err != nil {
		return 1
	}

	if opts.Version {
		fmt.Fprintf(os.Stderr, "%s/%s %s/%s\n", app, version, runtime.GOOS, runtime.GOARCH)

		return 0
	}

	if err := configen.Generate(&opts.Options, opts.Env...); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return 1
	}

	if opts.Watch {
		if err := configen.Watch(opts.Port, &opts.Options, opts.Env...); err != nil {
			fmt.Fprintln(os.Stderr, err)

			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(run(os.Args))
}
