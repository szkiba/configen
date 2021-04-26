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
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/jpillora/longestcommon"
)

// Watch start watching for changes of input files/directories and
// call Generate when change happened.
func Watch(port int, opts *Options, envs ...string) error {
	srv, err := newServer(port, opts, envs...)
	if err != nil {
		return err
	}

	return srv.run()
}

type server struct {
	watcher *fsnotify.Watcher
	dir     string
	opts    *Options
	envs    []string
	port    int
}

func newServer(port int, opts *Options, envs ...string) (*server, error) {
	srv := new(server)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	srv.watcher = watcher

	srv.opts = opts
	srv.envs = envs
	srv.port = port

	if err := srv.init(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (s *server) run() error {
	defer s.watcher.Close()

	done := make(chan bool)

	go s.watch()

	http.Handle("/", http.FileServer(http.Dir(s.dir)))

	listener, err := net.Listen("tcp", addr(s.port))
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Listening on http://%s\n", addr(listener.Addr().(*net.TCPAddr).Port))

	if err := http.Serve(listener, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	<-done

	return nil
}

func (s *server) onCreate(path string) {
	if _, err := os.Stat(path); err == nil {
		if err := s.watcher.Add(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (s *server) onModify() {
	fmt.Fprint(os.Stderr, "Change detected, generating ... ")

	if err := Generate(s.opts, s.envs...); err != nil {
		fmt.Fprintln(os.Stderr, "failed")
		fmt.Fprintln(os.Stderr, err)

		return
	}

	fmt.Fprintln(os.Stderr, "done")
}

func (s *server) watch() {
	for {
		select {
		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				s.onCreate(event.Name)
			}

			s.onModify()

		case err, ok := <-s.watcher.Errors:
			if !ok {
				return
			}

			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (s *server) init() error {
	outs := make([]string, 0, len(s.envs)+1)
	set := make(map[string]bool)

	outs = append(outs, s.opts.Output)

	for _, env := range s.envs {
		out, err := resolve(env, s.opts.Output)
		if err != nil {
			return err
		}

		outs = append(outs, out)

		if len(s.opts.Package) != 0 {
			if err := resolveToMap(env, []string{s.opts.Package}, set); err != nil {
				return err
			}
		}

		paths := make([]string, 0, len(s.opts.Templates)+len(s.opts.Raws)+len(s.opts.Schemas)+len(s.opts.Values))

		paths = append(append(append(append(paths, s.opts.Templates...),
			s.opts.Raws...),
			s.opts.Schemas...),
			s.opts.Values...)

		if err := resolveToMap(env, paths, set); err != nil {
			return err
		}
	}

	s.dir = strings.TrimSuffix(longestcommon.Prefix(outs), string([]byte{filepath.Separator}))

	for k := range set {
		if err := watchDeep(s.watcher, k); err != nil {
			return err
		}
	}

	return nil
}

func resolveToMap(env string, strs []string, set map[string]bool) error {
	for _, str := range strs {
		s, err := resolve(env, str)
		if err != nil {
			return err
		}

		set[s] = true
	}

	return nil
}

func watchDeep(watcher *fsnotify.Watcher, path string) error {
	if err := watcher.Add(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return watcher.Add(path)
	})
}

func addr(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}
