// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file is modified from os/exec/lp_plan.go

package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// ErrNotFound is the error resulting if a path search failed to find an executable file.
var ErrNotFound = errors.New("executable file not found in $path")

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}

// lookPath searches for an executable binary named file
// in the directories named by the path environment variable.
// If file begins with "/", "#", "./", or "../", it is tried
// directly and the path is not consulted.
// The result may be an absolute path or a path relative to the current directory.
func lookPath(file string) (string, error) {
	// skip the path lookup for these prefixes
	skip := []string{"/", "#", "./", "../"}

	for _, p := range skip {
		if strings.HasPrefix(file, p) {
			err := findExecutable(file)
			if err == nil {
				return file, nil
			}
			return "", &exec.Error{file, err}
		}
	}

	path := os.Getenv("path")
	for _, dir := range strings.Split(path, ":") {
		if err := findExecutable(dir + "/" + file); err == nil {
			return dir + "/" + file, nil
		}
	}
	return "", &exec.Error{file, ErrNotFound}
}
