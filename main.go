package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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

	// Parse flags
	flag.Parse()

	// Check for version flag
	if *versionFlag || command == "version" {
		fmt.Printf("Version: %s\n", version)
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

	fmt.Println("Files found:", len(filesToWork))

	// Print results if view flag is set
	if *viewFlag || *vFlag {
		printFiles(filesToWork)
	}

	applyFixes(filesToWork)
	// Apply fixes if fix flag is set
	if *fixFlag {
	}
}
