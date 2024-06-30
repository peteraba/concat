package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
)

func main() {
	workdir, dryRun := getArgs()

	fileParts := findParts(workdir)
	for target, parts := range fileParts {
		sort.Strings(parts)

		err := validateFileParts(target, parts)
		if err != nil {
			fmt.Println(err)
			continue
		}

		processFileParts(target, parts, dryRun)
	}
}

func getArgs() (string, bool) {
	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dryRun := false

	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-d" || os.Args[i] == "--dry-run" {
			dryRun = true
		}
	}

	if len(os.Args) > 1 {
		lastArg := os.Args[len(os.Args)-1]
		if lastArg != "-d" && lastArg != "--dry-run" {
			workdir = lastArg
		}
	}

	return workdir, dryRun
}

func validateFileParts(target string, parts []string) error {
	for i := 0; i < len(parts); i++ {
		expected := fmt.Sprintf("%s.%03d", target, i+1)
		if parts[i] != expected {
			return fmt.Errorf("wrong file name found. expected: '%s', found: '%s', file number: #%d", expected, parts[i], i+1)
		}
	}

	_, err := os.Stat(target)
	if err == nil {
		return fmt.Errorf("file '%s' already exists", target)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat file '%s': %w", target, err)
	}

	return nil
}

func processFileParts(target string, parts []string, dryRun bool) {
	if dryRun {
		return
	}

	err := writeFile(target, parts)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("file created:", target)
	}
}

func findParts(dir string) map[string][]string {
	result := make(map[string][]string)

	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("failed to read directory '%s': %v\n", dir, err)

		return nil
	}

	r, _ := regexp.Compile("\\.([0-9]{3})$")

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		match := r.FindString(name)
		if match == "" {
			continue
		}

		short := name[:len(name)-4]

		result[short] = append(result[short], name)
	}

	return result
}

func writeFile(target string, parts []string) error {
	f, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", target, err)
	}
	defer f.Close()

	for _, part := range parts {
		dat, err := os.ReadFile(part)
		if err != nil {
			panic(err)
		}

		f.Write(dat)
	}

	return nil
}
