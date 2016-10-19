// Copyright 2016 Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>s. All rights reserved.
//
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	shutil "github.com/termie/go-shutil"
	libsass "github.com/wellington/go-libsass"
)

const (
	slateRepo = "https://github.com/lord/slate.git"
)

var (
	logger *log.Logger = log.New(os.Stdout, "bundler: ", log.Ldate|log.Ltime|log.Lshortfile)

	staticSlateDirs = []string{
		"images",
		"fonts",
	}
)

func main() {

	var (
		slateSourceDir = flag.String("slate", "", "the path to the Slate source, if not set it will be cloned from "+slateRepo)
	)

	flag.Parse()

	pwd, err := os.Getwd()

	if err != nil {
		logger.Fatal(err)
	}

	bundler := newBundler(
		*slateSourceDir,
		filepath.Join(pwd, "static", "slate"))

	if err := bundler.init(); err != nil {
		logger.Fatal(err)
	}

	if err := bundler.fetchSlateIfNeeded(); err != nil {
		logger.Fatal(err)
	}

	if err := bundler.replaceSlateSourcesInTheme(); err != nil {
		logger.Fatal("Failed to move Slate sources: ", err)
	}

	if err := bundler.fixFontPaths(); err != nil {
		logger.Fatal("Failed to move Slate sources: ", err)
	}

	if err := bundler.createJSBundles(); err != nil {
		logger.Fatal("Failed to bundle JS: ", err)
	}

	if err := bundler.compileSassSources(); err != nil {
		logger.Fatal("Failed compile SASS stylesheets: ", err)
	}

}

type bundler struct {
	slateSource string
	slateTarget string
	logger      *log.Logger
}

func newBundler(slateSource, slateTarget string) *bundler {
	return &bundler{slateSource: slateSource, slateTarget: slateTarget, logger: logger}
}

func (b *bundler) init() error {
	if err := os.RemoveAll(b.slateTarget); err != nil {
		return err
	}

	if err := os.MkdirAll(b.slateTarget, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (b *bundler) fetchSlateIfNeeded() error {
	if b.slateSource != "" {
		b.logger.Println("Use existing Slate clone in", b.slateSource)
		return nil
	}

	b.logger.Println("Fetch Slate from", slateRepo)

	slateSource, err := ioutil.TempDir("", "hugo-slate")

	if err != nil {
		return fmt.Errorf("Failed to create tmpdir: %s", err)
	}

	if err := cloneSlateInto(slateSource); err != nil {
		return fmt.Errorf("Failed to clone Slate: %s", err)
	}

	b.slateSource = slateSource

	return nil
}

func (b *bundler) replaceSlateSourcesInTheme() error {
	for _, staticDir := range staticSlateDirs {
		b.logger.Println("Copy", staticDir)
		if err := shutil.CopyTree(filepath.Join(b.slateSource, "source", staticDir), filepath.Join(b.slateTarget, staticDir), nil); err != nil {
			return err
		}
	}
	return nil
}

// TODO(bep) this should be handled by the SASS compiler
func (b *bundler) fixFontPaths() error {
	base := filepath.Join(b.slateSource, "source", "stylesheets")
	for _, filename := range []string{"_icon-font.scss"} {
		fp := filepath.Join(base, filename)
		read, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}
		nc := bytes.Replace(read, []byte("font-url('"), []byte("font-url('../fonts/"), -1)

		err = ioutil.WriteFile(fp, nc, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bundler) compileSassSources() error {
	source := filepath.Join(b.slateSource, "source", "stylesheets")
	target := filepath.Join(b.slateTarget, "stylesheets")

	os.MkdirAll(target, os.ModePerm)

	fis, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), "_") {
			continue
		}

		targetName := strings.TrimSuffix(fi.Name(), ".scss")

		b.logger.Println("Compile", fi.Name(), "to", targetName)

		cssFile, err := os.Create(filepath.Join(target, targetName))
		if err != nil {
			return err
		}

		comp, err := libsass.New(cssFile, nil, libsass.Path(filepath.Join(source, fi.Name())))
		if err != nil {
			return err
		}

		if err := comp.Run(); err != nil {
			return err
		}

	}

	return nil
}

func (b *bundler) createJSBundles() error {
	src := filepath.Join(b.slateSource, "source", "javascripts")
	dst := filepath.Join(b.slateTarget, "javascripts")
	overrides := filepath.Join(b.slateTarget, "..", "..", "assets", "javascripts")
	jsB := newJSBundler(src, dst, overrides)
	return jsB.bundle()
}

func cloneSlateInto(dir string) error {
	logger.Println("Clone Slate into", dir, "...")

	cmd := exec.Command("git", "clone", slateRepo, dir)
	return cmd.Run()
}

type jsBundler struct {
	src string
	dst string

	overridesSrc string
	overrides    map[string][]byte

	logger *log.Logger

	// Per bundle
	seen map[string]bool
	buff bytes.Buffer
}

func newJSBundler(src, dst, overridesSrc string) *jsBundler {
	return &jsBundler{src: src, dst: dst, overridesSrc: overridesSrc, logger: logger, overrides: make(map[string][]byte)}
}

func (j *jsBundler) readOverrides() error {
	j.logger.Println("Looking for overrides in", j.overridesSrc)
	return filepath.Walk(j.overridesSrc, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		libPath := strings.TrimPrefix(path, j.overridesSrc)
		libPath = strings.TrimPrefix(libPath, string(filepath.Separator))
		j.logger.Println("Adding override:", libPath)

		libContent, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		j.overrides[libPath] = libContent

		return nil
	})

}

func (j *jsBundler) bundle() error {

	if err := j.readOverrides(); err != nil {
		return err
	}

	if err := os.MkdirAll(j.dst, os.ModePerm); err != nil {
		return err
	}

	fis, err := ioutil.ReadDir(j.src)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if !strings.HasSuffix(fi.Name(), ".js") {
			continue
		}

		filename := filepath.Join(j.src, fi.Name())
		if err := j.newBundle(filename); err != nil {
			return err
		}

		if err = ioutil.WriteFile(filepath.Join(j.dst, fi.Name()), j.buff.Bytes(), os.ModePerm); err != nil {
			return fmt.Errorf("Failed to write to destination: %s", err)
		}

	}

	return nil
}

func (j *jsBundler) newBundle(filename string) error {
	j.logger.Println("New bundle from ", filename)
	j.seen = make(map[string]bool)
	j.buff.Reset()

	j.buff.WriteString("\n\n// From bep's Poor Man's JS bundler // ----\n")

	return j.handleFile(filename)
}

func (j *jsBundler) handleFile(filename string) error {

	var (
		overridenContent = j.getOverrideIfFound(filename)
		currDir          = filepath.Dir(filename)
		libs             []string
	)

	if overridenContent == nil {
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		// TODO(bep) exclude the requires when writing to bundle
		libs = j.extractRequiredLibs(file)

		file.Close()
	} else {
		libs = j.extractRequiredLibs(bytes.NewReader(overridenContent))
	}

	for _, lib := range libs {
		if j.seen[lib] {
			continue
		}

		j.seen[lib] = true

		lib += ".js"
		libFilename := filepath.Join(currDir, lib)

		//j.logger.Println("Handle lib", lib)

		// Must write its dependencies first
		if err := j.handleFile(libFilename); err != nil {
			return err
		}

		var content []byte
		var err error

		content = j.getOverrideIfFound(libFilename)

		if content == nil {
			content, err = ioutil.ReadFile(libFilename)
			if err != nil {
				return err
			}
		}

		_, err = j.buff.Write(content)
		if err != nil {
			return err
		}

	}

	return nil
}

func (j *jsBundler) getOverrideIfFound(filename string) []byte {
	libPath := strings.TrimPrefix(filename, j.src)
	libPath = strings.TrimPrefix(libPath, string(filepath.Separator))
	return j.overrides[libPath]
}

func (j *jsBundler) extractRequiredLibs(r io.Reader) []string {
	const require = "//= require"
	scanner := bufio.NewScanner(r)
	var libs []string
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(t, require) {
			libs = append(libs, strings.TrimSpace(t[len(require):]))
		}
	}
	return libs
}
