package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	lines       int
	fileScanner *bufio.Scanner

	inputScanner = bufio.NewScanner(os.Stdin)
)

func main() {
	fmt.Print(":: Enter the directory path you are trying to count from (case sensitive):\n:: ")
	inputScanner.Scan()

	path := strings.TrimSpace(inputScanner.Text())

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println(err)
		return
	}

	fmt.Print(":: Would you like to remove empty lines? [y/N]: ")
	inputScanner.Scan()

	var remove bool

	switch sanitize(inputScanner.Text()) {
	case "y":
		remove = true
	case "n":
		remove = false
	case "":
		remove = false
	default:
		fmt.Println(":: Invalid option, not removing empty lines.")
		remove = false
	}

	fmt.Print(":: Enter the file extension(s) would you like to include, separated by commas. Leave blank to count all files:\n:: ")
	inputScanner.Scan()

	var extensions []string
	if strings.TrimSpace(inputScanner.Text()) != "" {
		extensions = strings.Split(inputScanner.Text(), ",")
		for i, ext := range extensions {
			extensions[i] = sanitize(ext)
		}
	}

	if !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		var count bool
		for _, ext := range extensions {
			if strings.HasSuffix(file.Name(), ext) {
				count = true
				break
			}
		}

		if !count {
			fmt.Println(":: The provided path lead to a file that does not contain any of the specified extensions.")
			return
		}

		fileScanner = bufio.NewScanner(file)
		for fileScanner.Scan() {
			if remove && strings.TrimSpace(fileScanner.Text()) == "" {
				continue
			}
			lines++
		}

		fmt.Println("This path did not lead to a directory, but the line count of this file is " + strconv.Itoa(lines) + ".")
	}

	err = filepath.Walk(path, func(p string, i fs.FileInfo, e error) error {
		if e != nil {
			return e
		}

		count := len(extensions) == 0
		for _, ext := range extensions {
			if strings.HasSuffix(i.Name(), ext) {
				count = true
				break
			}
		}

		if i.IsDir() || !count {
			return nil
		}

		file, err := os.Open(p)
		if err != nil {
			return err
		}

		fileScanner = bufio.NewScanner(file)
		for fileScanner.Scan() {
			if remove && strings.TrimSpace(fileScanner.Text()) == "" {
				continue
			}
			lines++
		}

		return nil
	})

	fmt.Println("The line count of all files in this directory is " + strconv.Itoa(lines) + ".")
}

func sanitize(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}
