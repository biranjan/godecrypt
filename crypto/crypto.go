package crpypto

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Trim(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s := strings.TrimSuffix(s, suffix)
		return s
	}
	return s
}

func GatherFiles(root string) []string {
	inputfiles := []string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic("Cannot walk in given path")
		}
		if !info.IsDir() {
			if strings.HasSuffix(path, ".c4gh") {
				inputfiles = append(inputfiles, path)
			}
		}
		return nil
	})
	if err != nil {
		panic("problem with file walk")
	}

	return inputfiles
}

func Decrypt(path string) (string, error) {
	log.Printf("Start processing %s\n", path)
	task := "decrypt"
	keypath := os.Getenv("key_path")
	log.Printf("key path %s", keypath)
	stdoutfile := Trim(path, ".c4gh")
	stdinfile := path
	subprocess := exec.Command("crypt4gh", task, "--sk", keypath)
	stdin, _ := subprocess.StdinPipe()
	fileout, err := os.Create(stdoutfile)
	if err != nil {
		return "", fmt.Errorf("error creating file %s, with %s", stdoutfile, err)
	}
	subprocess.Stdout = fileout
	stderr := subprocess.Stderr
	defer fileout.Close()
	filein, err := os.Open(stdinfile)
	if err != nil {
		return "", fmt.Errorf("error opening file %s, with %s", stdinfile, err)
	}
	defer filein.Close()
	go func() {
		io.Copy(stdin, filein)
		defer stdin.Close()
	}()

	err = subprocess.Start()
	if err != nil {
		return "", fmt.Errorf("error on job %s, with %s", path, err)
	}
	if err := subprocess.Wait(); err != nil {
		return "", fmt.Errorf("error on job %s, with %s", path, err)
	}
	log.Printf("std err %s ", stderr)

	if err != nil {
		return "", fmt.Errorf("error on job %s, with %s", path, err)
	} else {
		if err != nil {
			return "", fmt.Errorf("error on job %s, with %s", path, err)
		}

		return stdinfile, nil
	}

}
