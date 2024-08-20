package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func mergeFlags(long, short string) []string {
	if long != "" {
		return splitAndTrim(long)
	}
	return splitAndTrim(short)
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func findFiles2(files []string, fileRegex, ignoreRegex *regexp.Regexp, ignoreFiles map[string]interface{}) (map[string]interface{}, error) {
	filesMap := make(map[string]interface{})

	for _, file := range files {
		absPath, _ := filepath.Abs(file)
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", absPath)
		}
		if err := findMatchingFiles(absPath, fileRegex, ignoreRegex, filesMap, ignoreFiles); err != nil {
			return nil, fmt.Errorf("error processing path %s: %v", absPath, err)
		}
	}
	return filesMap, nil
}

func findMatchingFiles(path string, fileRegex, ignoreRegex *regexp.Regexp, oldFiles, ignoreFiles map[string]interface{}) error {
	return filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				fmt.Printf("Warning: Permission denied accessing %s\n", file)
				return nil
			}
			return fmt.Errorf("error accessing %s: %v", file, err)
		}

		if !info.IsDir() {
			basename := filepath.Base(file)
			if _, ok := ignoreFiles[file]; !ok {
				if fileRegex.MatchString(basename) && (ignoreRegex == nil || !ignoreRegex.MatchString(basename)) && oldFiles != nil {
					oldFiles[file] = struct{}{}
				}
			}
		}

		return nil
	})
}
func pathMapToArr(pathsMap map[string]interface{}) []string {
	arr := make([]string, 0, 30)
	for k, _ := range pathsMap {
		arr = append(arr, k)
	}
	return arr
}
func printFiles(files map[string]interface{}) {
	if len(files) == 0 {
		return
	}

	// Создаем слайс для хранения ключей (путей файлов)
	keys := make([]string, 0, len(files))
	for file := range files {
		keys = append(keys, file)
	}

	// Сортируем слайс
	sort.Strings(keys)

	fmt.Println("------------------------")
	for _, file := range keys {
		fmt.Println(file)
	}
	fmt.Println("------------------------")
}

func printUsage() {
	fmt.Println("Usage of go-padding:")
	fmt.Println("  go-padding --file <files> [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --file, -f            Comma-separated list of files or folders to process (required)")
	fmt.Println("  --ignore, -i          Comma-separated list of files or folders to ignore")
	fmt.Println("  --view, -v            Print the absolute paths of found files")
	fmt.Println("  --fix                 Make changes to the files")
	fmt.Println("  --pattern   		  Regex pattern for files to process (default: \\.go$)")
	fmt.Println("  --ignore-pattern	  Regex pattern for files to ignore")
	fmt.Println("  --version             Print the version of the program")
	fmt.Println("  --help                Print this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  go-padding --file folder1,folder2 --ignore folder/ignore")
	fmt.Println("  go-padding -f \"folder1, folder2/\" -i \"folder/ignore, folder2/ignore\"")
	fmt.Println("  go-padding --file folder1 --pattern \"\\.(go|txt)$\" --view")
	fmt.Println("  go-padding --file \"example, example, example/ignore\" --pattern \"(_test\\.go$|^filename_)\" --ignore-pattern \"_ignore\\.go$\" --view")
	fmt.Println("  go-padding --file \"example, example/userx_test.go\" --ignore-pattern \"_test\\.go|ignore\\.go$\" -v")
	fmt.Println("  go-padding --file \"pkg\"")
	fmt.Println("  go-padding --file \"pkg\" --ignore-pattern \"_test\\.go$\"")
	fmt.Println("  go-padding --file \"pkg\" --pattern \"_test\\.go$\"")
	fmt.Println("  go-padding --file folder1 --fix")
}
func applyFixes(files map[string]interface{}) {
	fmt.Println("Applying fixes to files:")
	for _, file := range files {
		fmt.Printf("Fixing file: %s\n", file)
		// Здесь должна быть логика внесения изменений в файл
		// Например:
		// applyFixToFile(file)
	}
}
