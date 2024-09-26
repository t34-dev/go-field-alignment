package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	version "github.com/t34-dev/go-field-alignment/v2"
)

// defaultFilePattern is the default regex pattern for files to process
const defaultFilePattern = `\.go$`

// main is the entry point of the program.
// It handles command-line arguments, processes files based on the provided flags,
// and applies necessary operations on the found files.
func main() {
	// Init logging.
	//
	// Error logging will go to stderr.
	log.SetFlags(0)

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
		log.Fatalf("Error compiling file pattern regex: %v\n", err)
	}

	var ignoreRegex *regexp.Regexp
	if ignorePattern != "" {
		ignoreRegex, err = regexp.Compile(ignorePattern)
		if err != nil {
			log.Fatalf("Error compiling ignore pattern regex: %v\n", err)
		}
	}

	ignoresMap, err := findFiles(ignores, fileRegex, ignoreRegex, nil)
	if err != nil {
		log.Fatalf("Cannot find files to ignore: %v\n", err)
	}
	filesToWork, err := findFiles(files, fileRegex, ignoreRegex, ignoresMap)
	if err != nil {
		log.Fatalf("Cannot find files to process: %v\n", err)
	}

	fmt.Printf("Files analyzed: %d\n-----------------\n", len(filesToWork))

	allFiles := make([]string, 0, len(filesToWork))
	for filePath := range filesToWork {
		allFiles = append(allFiles, filePath)
	}
	sort.Strings(allFiles)

	processingOpts := fileProcessingOptions{
		viewMode:  viewMode,
		fixMode:   fixMode,
		debugMode: debugMode,
	}

	var filesToFix []string
	for _, filePath := range allFiles {
		needFix, err := processFile(filePath, processingOpts)
		if err != nil {
			log.Fatalf("Cannot process file '%s': %v\n", filePath, err)
		}
		if needFix {
			filesToFix = append(filesToFix, filePath)
		}
	}
	if len(filesToFix) == 0 {
		return
	}

	if fixMode {
		fmt.Printf("-----------------\nApplied fixes to %d files\n", len(filesToFix))
	} else {
		fmt.Printf("-----------------\nFound files that need to be optimized:\n-- %s\n", strings.Join(filesToFix, "\n-- "))
		os.Exit(1)
	}
	fmt.Println()
}

// printUsage prints the usage information for the program.
func printUsage() {
	fmt.Println("Usage of gofield:")
	fmt.Println("  gofield --files <files> [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  --files, -f            Comma-separated list of files or folders to process (required)")
	fmt.Println("  --ignore, -i          Comma-separated list of files or folders to ignore")
	fmt.Println("  --view, -v            Print the absolute paths of found files")
	fmt.Println("  --fix                 Make changes to the files")
	fmt.Println("  --pattern   		  Regex pattern for files to process (default: \\.go$)")
	fmt.Println("  --ignore-pattern	  Regex pattern for files to ignore")
	fmt.Println("  --version             Print the version of the program")
	fmt.Println("  --help                Print this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  gofield --files folder1,folder2 --ignore folder/ignore")
	fmt.Println("  gofield -f \"folder1, folder2/\" -i \"folder/ignore, folder2/ignore\"")
	fmt.Println("  gofield --files folder1 --pattern \"\\.(go|txt)$\" --view")
	fmt.Println("  gofield --files \"example, example, example/ignore\" --pattern \"(_test\\.go$|^filename_)\" --ignore-pattern \"_ignore\\.go$\" --view")
	fmt.Println("  gofield --files \"example, example/userx_test.go\" --ignore-pattern \"_test\\.go|ignore\\.go$\" -v")
	fmt.Println("  gofield --files \"example\"")
	fmt.Println("  gofield --files \"example\" --ignore-pattern \"_test\\.go$\"")
	fmt.Println("  gofield --files \"example\" --pattern \"_test\\.go$\"")
	fmt.Println("  gofield --files example --fix")
}
