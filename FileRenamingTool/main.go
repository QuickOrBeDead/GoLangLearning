package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

func main() {
	root := "./files"
	// original := filepath.Join(root, "original")
	output := filepath.Join(root, "output")

	os.RemoveAll(output)
	os.Mkdir(output, 0655)
	// TODO: copy original files to output folder

	filepath.Walk(output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("path walk error", err)
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".txt" {
			renameFile(path)
		}

		return nil
	})
}

func renameFile(path string) {
	dir, filename := filepath.Split(path)
	regex, err := regexp.Compile(`(.*)_([0-9]+).txt`)
	if err != nil {
		fmt.Println("regex compile error: ", err.Error())
	}

	matches := regex.FindStringSubmatch(filename)
	if len(matches) > 2 {
		n, _ := strconv.Atoi(matches[2])
		total := int(math.Pow10(len(matches[2])))
		newFile := filepath.Join(dir, fmt.Sprintf("%s (%d of %d).txt", matches[1], n, total))
		err = os.Rename(path, newFile)
		if err != nil {
			fmt.Println("file rename error: ", err.Error())
			fmt.Printf("file '%s' is renamed to '%s'\n", path, newFile)
		}
	} else {
		fmt.Printf("file '%s' is not renamed\n", path)
	}
}
