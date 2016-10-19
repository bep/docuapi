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

type jsBundler struct {
	src string
	dst string
}

func (j *jsBundler) bundle() error {

	//	var bundle bytes.Buffer
	seen := make(map[string]bool)

	fis, err := ioutil.ReadDir(j.src)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if !strings.HasSuffix(fi.Name(), ".js") {
			continue
		}
		fmt.Println("Handle file", fi.Name())
		filename := filepath.Join(j.src, fi.Name())
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		libs := extractRequiredLibs(file)
		file.Close()
		for _, lib := range libs {
			if seen[lib] {
				continue
			}
			seen[lib] = true

			lib += ".js"
			libFilename := filepath.Join(j.src, lib)
			fmt.Println(">>>", libFilename)

		}

	}

	return nil
}

func extractRequiredLibs(r io.Reader) []string {
	const require = "//= require"
	scanner := bufio.NewScanner(r)
	var libs []string
	for scanner.Scan() {
		t := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(t, require) {
			libs = append(libs, t[len(require):])
		}
	}
	return libs
}

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

	/*if true {
		b := jsBundler{"/Users/bep/dev/clone/slate/source/javascripts", filepath.Join(slateTarget, "javascripts")}
		err := b.bundle()
		if err != nil {
			log.Fatal(err)
		}
		return
	}*/

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

	if err := fetchJavaScripts(filepath.Join(slateTarget, "javascripts")); err != nil {
		log.Fatal("Failed to fetch JS: ", err)
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
