package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	workdir, dryRun, force := getArgs(os.Args)

	fileParts := findParts(workdir)
	for target, parts := range fileParts {
		sort.Strings(parts)

		err := validateFileParts(target, parts)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = validateTarget(target)
		if err != nil {
			fmt.Println(err)
			if force {
				target = findNewTarget(target)
			} else {
				continue
			}
		}

		processFileParts(target, parts, dryRun, force)
	}
}

func getArgs(args []string) (string, bool, bool) {
	fmt.Println(args)

	workdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dryRun := false
	force := false

	for i := 1; i < len(args); i++ {
		if args[i] == "-f" || args[i] == "--force" {
			force = true
		} else if args[i] == "-d" || args[i] == "--dry-run" {
			dryRun = true
		} else if len(args[i]) > 2 && args[i][:2] == "--" {
			panic(fmt.Sprintf("unknown argument: '%s'", args[i]))
		}
	}

	if len(args) > 1 {
		lastArg := args[len(os.Args)-1]
		if lastArg != "-d" && lastArg != "--dry-run" && lastArg != "-f" && lastArg != "--force" {
			workdir = lastArg
		}
	}

	expectedArgCount := 2
	if dryRun {
		expectedArgCount++
	}
	if force {
		expectedArgCount++
	}
	if len(args) > expectedArgCount {
		panic("too many arguments")
	}

	return workdir, dryRun, force
}

func validateFileParts(target string, parts []string) error {
	for i := 0; i < len(parts); i++ {
		expected := fmt.Sprintf("%s.%03d", target, i+1)
		if parts[i] != expected {
			return fmt.Errorf("wrong file name found. expected: '%s', found: '%s', file number: #%d", expected, parts[i], i+1)
		}
	}

	return nil
}

func validateTarget(target string) error {
	_, err := os.Stat(target)
	if err == nil {
		return fmt.Errorf("file '%s' already exists", target)
	}
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to stat file '%s': %w", target, err)
	}

	return nil
}

func findNewTarget(target string) string {
	parts := strings.Split(target, ".")
	ext := ""
	targetBase := target

	if len(parts) > 1 {
		ext = "." + parts[len(parts)-1]
		targetBase = strings.Join(parts[:len(parts)-1], ".")
	}

	newTarget := ""
	for i := 1; i < 1000; i++ {
		newTarget = targetBase + fmt.Sprintf("-%03d", i) + ext

		_, err := os.Stat(newTarget)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			return newTarget
		}
	}

	panic(fmt.Errorf("failed to find new target file name, target: '%s'", target))
}

func processFileParts(target string, parts []string, dryRun, force bool) {
	if dryRun {
		return
	}

	err := writeFile(target, parts)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("file created:", target)
	}

	for _, part := range parts {
		err = os.Remove(part)
		if err != nil {
			fmt.Println("failed to remove file:", err)
		} else {
			fmt.Println("file removed:", part)
		}
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
