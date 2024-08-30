package main

import (
	"flag"
	"fmt"
	version "github.com/t34-dev/go-field-alignment/v2"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

// defaultFilePattern is the default regex pattern for files to process
const defaultFilePattern = `\.go$`

// main is the entry point of the program.
// It handles command-line arguments, processes files based on the provided flags,
// and applies necessary operations on the found files.
func main() {
	command := ""
	if len(os.Args) == 2 {
		command = strings.TrimSpace(os.Args[1])
	}

	// Define flags
	filesFlag := flag.String("files", "", "Comma-separated list of files or folders to process")
	fFlag := flag.String("f", "", "Short form of --files")
	ignoreFlag := flag.String("ignore", "", "Comma-separated list of files or folders to ignore")
	iFlag := flag.String("i", "", "Short form of --ignore")
	viewFlag := flag.Bool("view", false, "Print the absolute paths of found files")
	vFlag := flag.Bool("v", false, "Short form of --view")
	fixFlag := flag.Bool("fix", false, "Make changes to the files")
	filePatternFlag := flag.String("pattern", "", "Regex pattern for files to process")
	ignorePatternFlag := flag.String("ignore-pattern", "", "Regex pattern for files to ignore")
	versionFlag := flag.Bool("version", false, "Print the version of the program")
	helpFlag := flag.Bool("help", false, "Print usage information")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")

	// Parse flags
	flag.Parse()

	// Check for version flag
	if *versionFlag || command == "version" {
		fmt.Printf("Version: %s\n", version.Version)
		return
	}

	// Check for help flag or missing required flags
	if *helpFlag || command == "help" || (*filesFlag == "" && *fFlag == "") {
		printUsage()
		return
	}

	// Merge short and long form flags
	files := mergeFlags(*filesFlag, *fFlag)
	ignores := mergeFlags(*ignoreFlag, *iFlag)
	filePattern := *filePatternFlag
	ignorePattern := *ignorePatternFlag
	debugMode := *debugFlag
	fixMode := *fixFlag
	viewMode := *viewFlag || *vFlag

	// Ensure filePattern is not empty
	if filePattern == "" {
		filePattern = defaultFilePattern
	}

	// Compile regex patterns
	fileRegex, err := regexp.Compile(filePattern)
	if err != nil {
		fmt.Printf("Error compiling file pattern regex: %v\n", err)
		os.Exit(1)
	}

	var ignoreRegex *regexp.Regexp
	if ignorePattern != "" {
		ignoreRegex, err = regexp.Compile(ignorePattern)
		if err != nil {
			fmt.Printf("Error compiling ignore pattern regex: %v\n", err)
			os.Exit(1)
		}
	}

	ignoresMap, err := findFiles(ignores, fileRegex, ignoreRegex, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	filesToWork, err := findFiles(files, fileRegex, ignoreRegex, ignoresMap)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Files analyzed: %d\n-----------------\n", len(filesToWork))

	allFiles := make([]string, 0, len(filesToWork))
	for filePath := range filesToWork {
		allFiles = append(allFiles, filePath)
	}
	sort.Strings(allFiles)

	codeExit := 0
	for _, filePath := range allFiles {
		// read files
		openFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalln(err)
		}
		structures, mapStructures, err := Parse(openFile)
		if err != nil {
			log.Fatalln(err)
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
		needFix := false
		for _, structure := range structures {
			if structure.MetaData.BeforeSize > structure.MetaData.AfterSize {
				codeExit += 1
				needFix = true
				break
			}
		}
		if viewMode || needFix {
			fmt.Printf("%s\n", filePath)
		}
		for idx, structure := range structures {
			if structure.MetaData.BeforeSize > structure.MetaData.AfterSize {
				alert := fmt.Sprintf("can free %d bytes", structure.MetaData.BeforeSize-structure.MetaData.AfterSize)
				if fixMode {
					alert = "Fixed"
				}
				fmt.Printf("%s%-15s %d(b) -> %d(b) %s!\n", strings.Repeat(" ", 3),
					structure.Name,
					structure.MetaData.BeforeSize,
					structure.MetaData.AfterSize,
					alert)
				if debugMode {
					oldStructure, ok := oldStructuresMapper[structure.Path]
					if ok {
						fmt.Printf("%s%-20s\n", strings.Repeat(" ", 9), "------------------------------------------ [BEFORE]")
						testPrintStructure(oldStructure, 9)
						fmt.Printf("%s%-20s\n", strings.Repeat(" ", 9), "------------------------------------------ [AFTER]")
						testPrintStructure(structure, 9)
					}
				}
				if idx != len(structures)-1 && debugMode {
					fmt.Println()
				}
			} else {
				if viewMode {
					fmt.Printf("%s%-15s ✓\n", strings.Repeat(" ", 3), structure.Name)
				}
			}
		}
		if viewMode && len(structures) > 0 {
			fmt.Println()
		}

		// FIX
		if !fixMode {
			continue
		}
		err = renderTextStructures(structures)
		if err != nil {
			log.Fatal(err)
		}
		// Replace content
		resultFile, err := Replacer(openFile, structures)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(filePath, resultFile, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	if codeExit > 0 {
		if !fixMode {
			fmt.Printf("-----------------\nFound files: %d. That need to be optimized.\n", codeExit)
			os.Exit(1)
		} else {
			fmt.Printf("-----------------\nApplying fixes to files: %d\n", codeExit)
		}
		fmt.Println()
	}
}
