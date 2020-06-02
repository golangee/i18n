// Copyright 2020 Torben Schinke
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ModRootDir returns the root directory of current module. If the current working directory is not a module
// returns an error.
func ModRootDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	root := cwd
	for {
		stat, err := os.Stat(filepath.Join(root, "go.mod"))
		if err == nil && stat.Mode().IsRegular() {
			return root, nil
		}
		root = filepath.Dir(root)
		if root == "/" || root == "." {
			return "", fmt.Errorf("%s is not withing a go module", cwd)
		}
	}
}

// ListFiles simply returns all absolute file names of regular files ignoring any errors.
func (p *Package) ListFiles() []string {
	var res []string

	files, err := ioutil.ReadDir(p.Dir)

	if err != nil {
		return res
	}

	for _, f := range files {
		if f.Mode().IsRegular() {
			res = append(res, filepath.Join(p.Dir, f.Name()))
		}
	}

	return res
}

func (p *Package) String() string {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// Package contains information similar to 'go list -json'
type Package struct {
	Dir      string
	Name     string
	Packages []*Package
}

// listRecursive is like GoList but scans through all available folders. The first encountered error is returned.
// However "cannot find module for path" errors are silently ignored, because they occur with empty or non-go folders.
func listRecursive(parent *Package) error {
	files, err := ioutil.ReadDir(parent.Dir)
	if err != nil {
		return fmt.Errorf("failed to read dir %s: %w", parent.Dir, err)
	}

	for _, file := range files {
		// ignore hidden
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		if file.IsDir() {
			childPath := filepath.Join(parent.Dir, file.Name())
			childPkg, err := GoList(childPath, true)

			if err != nil {
				return err
			}
			if childPkg != nil {
				parent.Packages = append(parent.Packages, childPkg)
			}
		}
	}

	return nil
}

// GoList is like 'go list -json' within the given directory, but faster.
func GoList(dir string, recursive bool) (*Package, error) {

	pkgName, err := guestimatePackageName(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to guess packagename from %s: %w", dir, err)
	}
	p := &Package{}
	p.Name = pkgName
	p.Dir = dir

	if recursive {
		err = listRecursive(p)
		if err != nil {
			return p, fmt.Errorf("failed to list packages: %w", err)
		}
	}

	return p, nil
}

func guestimatePackageName(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), ".go") {
			buf := make([]byte, 1024)
			file, err := os.Open(filepath.Join(dir, f.Name()))
			if err != nil {
				return "", err
			}
			n, err := file.Read(buf)
			_ = file.Close()
			if err != nil {
				return "", err
			}
			buf = buf[:n]
			lines := strings.Split(string(buf), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "package ") {
					return strings.TrimSpace(line[len("package "):]), nil
				}
			}
		}
	}

	return filepath.Base(dir), nil
}
