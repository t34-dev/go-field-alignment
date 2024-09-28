package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// mergeFlags combines the long and short form flags, prioritizing the long form if present.
func mergeFlags(long, short string) []string {
	if long != "" {
		return splitAndTrim(long)
	}
	return splitAndTrim(short)
}

// splitAndTrim splits a string by commas and trims whitespace from each part.
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

// findFiles searches for files matching the given regex patterns and ignoring specified files.
func findFiles(files []string, fileRegex, ignoreRegex *regexp.Regexp, ignoreFiles map[string]interface{}) (map[string]interface{}, error) {
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

// findMatchingFiles walks through the directory structure and finds files matching the given patterns.
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

// fileProcessingOptions is a set of options which define how file gets processed.
type fileProcessingOptions struct {
	viewMode  bool
	fixMode   bool
	debugMode bool
}

// processFile processes a file located at the specified path.
//
// Returns true if the file can be optimized (needs fix), false otherwise.
func processFile(path string, opts fileProcessingOptions) (needFix bool, err error) {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("cannot read file: %w", err)
	}
	structures, mapStructures, err := Parse(fileData)
	if err != nil {
		return false, fmt.Errorf("cannot parse file: %w", err)
	}

	calculateStructures(structures, true)

	oldStructures := make([]*Structure, 0, len(structures))
	for _, structure := range structures {
		copied := deepCopy(structure)
		oldStructures = append(oldStructures, copied)
	}
	oldStructuresMapper := createMapper(oldStructures)

	optimizeMapperStructures(mapStructures)
	calculateStructures(structures, false)

	for _, structure := range structures {
		if structure.MetaData.BeforeSize > structure.MetaData.AfterSize {
			needFix = true
			break
		}
	}
	if opts.viewMode || needFix {
		fmt.Printf("%s\n", path)
	}
	for idx, structure := range structures {
		if structure.MetaData.BeforeSize > structure.MetaData.AfterSize {
			alert := fmt.Sprintf("can free %d bytes", structure.MetaData.BeforeSize-structure.MetaData.AfterSize)
			if opts.fixMode {
				alert = "Fixed"
			}
			fmt.Printf(
				"%s%-15s %d(b) -> %d(b) %s!\n",
				strings.Repeat(" ", 3),
				structure.Name,
				structure.MetaData.BeforeSize,
				structure.MetaData.AfterSize,
				alert,
			)
			if opts.debugMode {
				oldStructure, ok := oldStructuresMapper[structure.Path]
				if ok {
					fmt.Printf("%s%-20s\n", strings.Repeat(" ", 9), "------------------------------------------ [BEFORE]")
					testPrintStructure(oldStructure, 9)
					fmt.Printf("%s%-20s\n", strings.Repeat(" ", 9), "------------------------------------------ [AFTER]")
					testPrintStructure(structure, 9)
				}
			}
			if idx != len(structures)-1 && opts.debugMode {
				fmt.Println()
			}
		} else {
			if opts.viewMode {
				fmt.Printf("%s%-15s ✓\n", strings.Repeat(" ", 3), structure.Name)
			}
		}
	}
	if opts.viewMode && len(structures) > 0 {
		fmt.Println()
	}

	if !opts.fixMode || !needFix {
		// If "fix" has not been requested or there's nothing to fix, exit
		return needFix, nil
	}

	// FIX
	renderTextStructures(structures)

	// Apply replacements
	resultData, err := Replacer(fileData, structures)
	if err != nil {
		return needFix, fmt.Errorf("cannot replace content in file: %w", err)
	}

	// Format results.
	//
	// They need to be formatted after all replacements have been applied
	formatted, err := format.Source(resultData)
	if err != nil {
		return needFix, fmt.Errorf("cannot format result content: %w", err)
	}

	// Write results
	err = os.WriteFile(path, formatted, 0644)
	if err != nil {
		return needFix, fmt.Errorf("cannot write results to file: %w", err)
	}
	return needFix, nil
}

// normalizeLineEndings converts all line endings to LF
func normalizeLineEndings(data []byte) []byte {
	return bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
}
