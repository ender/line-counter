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
	fmt.Print("Enter the absolute path to the folder you are trying to count lines from (path is case sensitive): ")
	inputScanner.Scan()

	path := strings.TrimSpace(inputScanner.Text())

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("This file or directory is nonexistent.")
		return
	}

	fmt.Print("Would you like to remove empty lines? (y/N) ")
	inputScanner.Scan()

	var remove bool

	switch sanitize(inputScanner.Text()) {
	case "y":
		remove = true
	case "n":
		remove = false
	default:
		fmt.Println("Invalid option, not removing empty lines.")
		remove = false
	}

	fmt.Print("What is the file extension of the files you would like to count? Leave blank to count all files. ")
	inputScanner.Scan()

	var extension string
	if sanitize(inputScanner.Text()) != "" {
		extension = sanitize(inputScanner.Text())
	}

	if !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
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

		if i.IsDir() || (extension != "" && !strings.HasSuffix(i.Name(), extension)) {
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
