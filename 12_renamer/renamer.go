package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type foundFile struct {
	path string
	info fs.FileInfo
}

func main() {

	var opts struct {
		CurrentDir         string `short:"d" long:"directory" description:"Directory to find files to rename" required:"true"`
		FindPatternRegex   string `short:"f" long:"find-pattern-regex" description:"Pattern regex to find files to replace" required:"true"`
		RenamePatternRegex string `short:"r" long:"rename-pattern-regex" description:"Pattern regex to rename files. Use <N> for current file and <T> for all files" required:"true"`
		DryRun             bool   `long:"dry-run" description:"Mock the renaming of file" required:"false"`
	}
	args, err := flags.Parse(&opts)
	// opts.CurrentDir := "sample"
	// opts.FindPatternRegex := "birthday_[0-9][0-9][0-9].txt"
	// opts.RenamePattern  := "birthday_(<N> out of <T>).txt"
	if err != nil {
		os.Exit(-1)
	}
	if len(args) != 0 {
		fmt.Println("Too many arguments:", args)
		os.Exit(-1)
	}
	re := regexp.MustCompile(opts.FindPatternRegex)
	foundFiles := findFilesToRename(opts.CurrentDir, re)
	fmt.Printf("Found %d file(s)\n", len(foundFiles))
	renameFiles(foundFiles, re, opts.RenamePatternRegex, opts.DryRun)
}

func findFilesToRename(currentDir string, re *regexp.Regexp) []foundFile {
	foundFiles := make([]foundFile, 0)
	filepath.Walk(currentDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matchesPattern(info, re) {
			foundFiles = append(foundFiles, foundFile{path: path, info: info})
		}
		return nil
	})
	return foundFiles
}

func matchesPattern(fileInfo fs.FileInfo, re *regexp.Regexp) bool {
	return re.MatchString(fileInfo.Name())
}

func renameFiles(foundFiles []foundFile, re *regexp.Regexp, renamePattern string, dryRun bool) {
	total := len(foundFiles)
	for i, file := range foundFiles {

		newName := re.ReplaceAllString(file.info.Name(), renamePattern)
		newName = strings.Replace(newName, "<N>", strconv.Itoa(i+1), -1)
		newName = strings.Replace(newName, "<T>", strconv.Itoa(total), -1)
		newPath := filepath.Join(filepath.Dir(file.path), newName)
		err := moveFilePath(file.info, file.path, newPath, dryRun)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func moveFilePath(info fs.FileInfo, oldPath, newPath string, dryRun bool) error {
	dryRunMessage := ""
	if dryRun {
		dryRunMessage = "[DRY-RUN]"
	}
	fmt.Printf("renaming file path %s to %s...%s\n", oldPath, newPath, dryRunMessage)
	if dryRun {
		return nil
	}
	return os.Rename(oldPath, newPath)
}
