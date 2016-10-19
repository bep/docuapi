package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	libsass "github.com/wellington/go-libsass"
)

func main() {

	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	// Store Slate statics below /static/slate
	slateTarget := filepath.Join(pwd, "static", "slate")

	if err := os.RemoveAll(slateTarget); err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(slateTarget, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update Slate from source ...")

	slateSource, err := ioutil.TempDir("", "hugo-slate")

	if err != nil {
		log.Fatal("Failed to create tmpdir: ", err)
	}

	if err := cloneSlateInto(slateSource); err != nil {
		log.Fatal("Failed to clone Slate: ", err)
	}

	if err := replaceSlateSourcesInTheme(slateSource, slateTarget); err != nil {
		log.Fatal("Failed to move Slate sources: ", err)
	}

	if err := fixFontPaths(filepath.Join(filepath.Join(slateSource, "source", "stylesheets"))); err != nil {
		log.Fatal("Failed to move Slate sources: ", err)
	}

	if err := createJSBundles(filepath.Join(slateSource, "source", "javascripts"), filepath.Join(slateTarget, "javascripts")); err != nil {
		log.Fatal("Failed to bundle JS: ", err)
	}

	if err := compileSassSources(filepath.Join(slateSource, "source", "stylesheets"), filepath.Join(slateTarget, "stylesheets")); err != nil {
		log.Fatal("Failed compile SASS stylesheets: ", err)
	}

}

func cloneSlateInto(dir string) error {
	fmt.Println("Clone Slate into", dir, "...")
	repo := "https://github.com/lord/slate.git"

	cmd := exec.Command("git", "clone", repo, dir)
	return cmd.Run()
}

var staticSlateDirs = []string{
	"images",
	"fonts",
}

func createJSBundles(src, dst string) error {
	b := newJSBUndler(src, dst)
	return b.bundle()
}

// TODO(bep) create JS bundles?
func fetchJavaScripts(target string) error {
	os.MkdirAll(target, os.ModePerm)
	for _, resource := range []string{"all.js", "all_nosearch.js"} {
		url := "https://lord.github.io/slate/javascripts/" + resource
		filepath := filepath.Join(target, resource)
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}

		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}

		resp.Body.Close()
		out.Close()

	}

	return nil

}

func replaceSlateSourcesInTheme(source, target string) error {
	for _, staticDir := range staticSlateDirs {
		fmt.Println("Move", staticDir)
		if err := os.Rename(filepath.Join(source, "source", staticDir), filepath.Join(target, staticDir)); err != nil {
			return err
		}
	}
	return nil
}

func compileSassSources(source, target string) error {
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
		fmt.Println("Compile", fi.Name(), "to", targetName)

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

func fixFontPaths(base string) error {
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

type jsBundler struct {
	src string
	dst string

	// Per bundle
	seen map[string]bool
	buff bytes.Buffer
}

func newJSBUndler(src, dst string) *jsBundler {
	return &jsBundler{src: src, dst: dst}
}

func (j *jsBundler) bundle() error {

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
	fmt.Println("New bundle from ", filename)
	j.seen = make(map[string]bool)
	j.buff.Reset()

	return j.handleFile(filename)
}

func (j *jsBundler) handleFile(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	libs := j.extractRequiredLibs(file)
	currDir := filepath.Dir(filename)
	file.Close()
	for _, lib := range libs {
		if j.seen[lib] {
			continue
		}
		j.seen[lib] = true

		lib += ".js"
		libFilename := filepath.Join(currDir, lib)

		fmt.Println("Handle lib", libFilename)

		// Must write the dependencies first.
		if err := j.handleFile(libFilename); err != nil {
			return err
		}

		// Write this to the bundle and continue
		content, err := ioutil.ReadFile(libFilename)
		if err != nil {
			return err
		}

		j.buff.WriteString(fmt.Sprintf("\n\n// bep's Poor Man's JS bundler: %s\n// ----\n", lib))
		_, err = j.buff.Write(content)
		if err != nil {
			return err
		}

	}

	return nil
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
