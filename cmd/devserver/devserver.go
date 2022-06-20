package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

const sourceURL = "http://opendata.ndw.nu/"
const cacheDir = "cache"

func placeFile(url, localpath string) error {
	var resp *http.Response
	var err error
	if resp, err = http.Get(url); err != nil {
		return err
	}

	var file []byte
	if file, err = io.ReadAll(resp.Body); err != nil {
		return err
	}

	return os.WriteFile(localpath, file, os.ModeType)
}

func ensureFiles() error {
	fileNames := []string{"DRIPS.xml.gz", "LocatietabelDRIPS.xml.gz"}

	if err := os.MkdirAll(cacheDir, os.ModeType); err != nil {
		return err
	}

	for _, f := range fileNames {
		filePath := path.Join(cacheDir, f)

		if _, err := os.Stat(filePath); err == nil {
			continue // File already exists
		}

		if err := placeFile(sourceURL+f, filePath); err != nil {
			return err
		}
	}

	return nil
}

func serve() error {
	fs := http.FileServer(http.Dir(cacheDir))
	path, err := filepath.Abs(cacheDir)
	if err != nil {
		fmt.Printf("Error parsing cacheDir")
		return err
	}

	fmt.Println("Serving dev data from ", path)
	return http.ListenAndServe("localhost:8001", fs)
}

func main() {
	if err := ensureFiles(); err != nil {
		fmt.Println(err)
		return
	}

	if err := serve(); err != nil {
		fmt.Println(err)
		return
	}
}
