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
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/russross/blackfriday/v2"
	"github.com/skratchdot/open-golang/open"
)

type markdownHandler struct {
	dir http.Dir
	h   http.Handler
}

type mdData struct {
	Title string
	Body  template.HTML
}

func addr(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}

func openReadme(o *options) error {
	dir := http.Dir(".")
	fs := http.FileServer(dir)
	h := &markdownHandler{dir: dir, h: fs}
	http.Handle("/", h)

	listener, err := net.Listen("tcp", addr(o.Port))
	if err != nil {
		return err
	}

	loc := fmt.Sprintf("http://%s/README.md", addr(listener.Addr().(*net.TCPAddr).Port))

	fmt.Fprintf(os.Stderr, "Opening %s\n", loc)

	done := make(chan bool)

	go func() {
		if err := http.Serve(listener, nil); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	open.Run(loc) // nolint

	<-done

	return nil
}

func (m *markdownHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.URL.Path, ".md") {
		m.h.ServeHTTP(rw, req)

		return
	}

	var pathErr *os.PathError

	md, err := ioutil.ReadFile(string(m.dir) + req.URL.Path)
	if errors.As(err, &pathErr) {
		http.Error(rw, http.StatusText(http.StatusNotFound)+": "+req.URL.Path, http.StatusNotFound)

		return
	}

	if err != nil {
		http.Error(rw, "Internal Server Error: "+err.Error(), 500)

		return
	}

	output := blackfriday.Run(md, blackfriday.WithExtensions(blackfriday.CommonExtensions|blackfriday.AutoHeadingIDs))

	rw.Header().Set("Content-Type", "text/html")

	data := &mdData{Title: req.URL.Path[1:], Body: template.HTML(string(output))} // nolint
	if err := mdTemplate.Execute(rw, data); err != nil {
		http.Error(rw, "Internal Server Error: "+err.Error(), 500)
	}
}

// nolint:lll
var mdTemplate = template.Must(template.New("md").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <title>{{ .Title }}</title>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"/>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css" crossorigin="anonymous" />    
  </head>
  <body>
    <div class="markdown-body">
    {{ .Body }}
    </div>
    <style>
      .markdown-body {
        max-width: 1280px;
        margin: 0 auto;
        padding: 1rem;
      }
    </style>
  </body>
</html>
`))
