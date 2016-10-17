package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	"javascripts",
	"images",
	"fonts",
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
