package main

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var fileNameRegex *regexp.Regexp

func main() {
	var err error
	fileNameRegex, err = regexp.Compile(`(.*)_([0-9]+).txt`)
	if err != nil {
		fmt.Println("regex compile error: ", err.Error())
	}

	root := "./files"
	original := filepath.Join(root, "original")
	output := filepath.Join(root, "output")

	os.RemoveAll(output)
	os.Mkdir(output, 0655)
	copyDirContent(original, output)

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

	matches := fileNameRegex.FindStringSubmatch(filename)
	if len(matches) > 2 {
		n, _ := strconv.Atoi(matches[2])
		total := int(math.Pow10(len(matches[2])))
		newFile := filepath.Join(dir, fmt.Sprintf("%s (%d of %d).txt", matches[1], n, total))
		err := os.Rename(path, newFile)
		if err != nil {
			fmt.Println("file rename error: ", err.Error())
			fmt.Printf("file '%s' is renamed to '%s'\n", path, newFile)
		}
	} else {
		fmt.Printf("file '%s' is not renamed\n", path)
	}
}

func copyDirContent(src string, dest string) error {
	if err := checkIsDir(src); err != nil {
		return err
	}

	if err := checkIsDir(dest); err != nil {
		return err
	}

	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if src == path {
			return nil
		}

		if info.IsDir() {
			folder := strings.TrimPrefix(path, src)
			os.Mkdir(filepath.Join(dest, folder), 0655)
		} else if info.Mode().IsRegular() {
			srcDir, _ := filepath.Split(path)
			folder := strings.TrimPrefix(srcDir, src)
			copyFile(path, filepath.Join(dest, folder, info.Name()))
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func checkIsDir(path string) error {
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}

	return nil
}

func copyFile(src string, dest string) error {
	srcStats, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !srcStats.Mode().IsRegular() {
		return fmt.Errorf("%s file is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()

	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("destination file %s already exists", dest)
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destFile.Close()

	buf := make([]byte, os.Getpagesize()*32)
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = destFile.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return err
}
