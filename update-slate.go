package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	pwd, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Update Slate from source ...")

	slateDir, err := ioutil.TempDir("", "hugo-slate")

	if err != nil {
		log.Fatal("Failed to create tmpdir: ", err)
	}

	if err := cloneSlateInto(slateDir); err != nil {
		log.Fatal("Failed to clone Slate: ", err)
	}

	if err := replaceSlateSourcesInTheme(slateDir, pwd); err != nil {
		log.Fatal("Failed to move Slate sources: ", err)
	}
}

func cloneSlateInto(dir string) error {
	fmt.Println("Clone Slate into", dir, "...")
	repo := "https://github.com/lord/slate.git"

	cmd := exec.Command("git", "clone", repo, dir)
	return cmd.Run()
}

var staticSlateDirs = []string{
	"stylesheets",
	"javascripts",
	"images",
	"fonts",
}

func replaceSlateSourcesInTheme(source, target string) error {
	for _, staticDir := range staticSlateDirs {
		fmt.Println("Move", staticDir)
		if err := os.Rename(filepath.Join(source, "source", staticDir), filepath.Join(target, "static", staticDir)); err != nil {
			return err
		}
	}
	return nil
}
