package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// TestGoPadding is the main test function that runs all subtests for the gopad program.
func TestGoPadding(t *testing.T) {
	// Setup: Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "gopad-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	createTestFiles(t, tempDir)

	// Run tests
	t.Run("BasicUsage", func(t *testing.T) { testBasicUsage(t, tempDir) })
	t.Run("IgnoreFiles", func(t *testing.T) { testIgnoreFiles(t, tempDir) })
	t.Run("ViewFiles", func(t *testing.T) { testViewFiles(t, tempDir) })
	t.Run("ApplyFixes", func(t *testing.T) { testApplyFixes(t, tempDir) })
	t.Run("FilePatterns", func(t *testing.T) { testFilePatterns(t, tempDir) })
	t.Run("ErrorHandling", func(t *testing.T) { testErrorHandling(t, tempDir) })
	t.Run("UtilityFunctions", func(t *testing.T) { testUtilityFunctions(t) })
	t.Run("FlagCombinations", func(t *testing.T) { testFlagCombinations(t, tempDir) })
	t.Run("PathsWithSpaces", func(t *testing.T) { testPathsWithSpaces(t, tempDir) })
	t.Run("RecursiveTraversal", func(t *testing.T) { testRecursiveTraversal(t, tempDir) })
	// Temporarily disable SymbolicLinks test on Windows
	if runtime.GOOS != "windows" {
		t.Run("SymbolicLinks", func(t *testing.T) { testSymbolicLinks(t, tempDir) })
	}
	t.Run("AccessRights", func(t *testing.T) { testAccessRights(t, tempDir) })
}

// createTestFiles creates a set of test files in the specified directory.
func createTestFiles(t *testing.T, dir string) {
	files := []string{"main.go", "utils.go", "test.go", "ignore_me.go"}
	for _, file := range files {
		err := os.WriteFile(filepath.Join(dir, file), []byte("package main"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create a subdirectory
	subDir := filepath.Join(dir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	err := os.WriteFile(filepath.Join(subDir, "subfile.go"), []byte("package sub"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file in subdirectory: %v", err)
	}
}

// runCommand executes the gopad program with the given arguments and returns the output.
func runCommand(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// testBasicUsage tests the basic usage of the gopad program.
func testBasicUsage(t *testing.T, dir string) {
	output, err := runCommand("--files", filepath.Join(dir, "main.go"))
	if err != nil {
		t.Errorf("Basic usage failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 1") {
		t.Errorf("Unexpected output for basic usage: %s", output)
	}

	output, err = runCommand("-f", fmt.Sprintf("%s,%s", filepath.Join(dir, "main.go"), filepath.Join(dir, "utils.go")))
	if err != nil {
		t.Errorf("Basic usage with short flag failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 2") {
		t.Errorf("Unexpected output for basic usage with short flag: %s", output)
	}
}

// testIgnoreFiles tests the file ignoring functionality of gopad.
func testIgnoreFiles(t *testing.T, dir string) {
	output, err := runCommand("--files", dir, "--ignore", filepath.Join(dir, "test.go")+","+filepath.Join(dir, "ignore_me.go"))
	if err != nil {
		t.Errorf("Ignore files failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 3") {
		t.Errorf("Unexpected output for ignore files: %s", output)
	}

	output, err = runCommand("-f", dir, "-i", fmt.Sprintf("%s,%s", filepath.Join(dir, "test.go"), filepath.Join(dir, "ignore_me.go")))
	if err != nil {
		t.Errorf("Ignore files with short flag failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 3") {
		t.Errorf("Unexpected output for ignore files with short flag: %s", output)
	}
}

// testViewFiles tests the file viewing functionality of gopad.
func testViewFiles(t *testing.T, dir string) {
	output, err := runCommand("--files", dir, "--view")
	if err != nil {
		t.Errorf("View files failed: %v", err)
	}
	if !strings.Contains(output, filepath.Join(dir, "main.go")) {
		t.Errorf("View files did not show expected file: %s", output)
	}

	output, err = runCommand("-f", dir, "-v")
	if err != nil {
		t.Errorf("View files with short flag failed: %v", err)
	}
	if !strings.Contains(output, filepath.Join(dir, "utils.go")) {
		t.Errorf("View files with short flag did not show expected file: %s", output)
	}
}

// testApplyFixes tests the fix applying functionality of gopad.
func testApplyFixes(t *testing.T, dir string) {
	filePath := filepath.Join(dir, "main.go")

	// Сохраняем исходное содержимое файла
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read original file content: %v", err)
	}
	// Функция для восстановления исходного содержимого
	defer func() {
		err := os.WriteFile(filePath, originalContent, 0644)
		if err != nil {
			t.Errorf("Failed to restore original file content: %v", err)
		}
	}()

	// Добавляем тестовую структуру с невыровненными полями
	testStructure := `package main

type TestStruct struct {
    a bool
    b int64
    c bool
}
`
	err = os.WriteFile(filePath, []byte(testStructure), 0644)
	if err != nil {
		t.Fatalf("Failed to write test structure to file: %v", err)
	}

	output, err := runCommand("--files", filePath, "--fix")
	if err != nil {
		t.Errorf("Apply fixes failed: %v", err)
	}
	if !strings.Contains(output, "Applying fixes to structures:") {
		t.Errorf("Unexpected output for apply fixes: %s", output)
	}
	// Проверяем, что структура была оптимизирована
	optimizedContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read optimized file content: %v", err)
	}

	expectedOptimizedStructure := `package main

type TestStruct struct {
	b int64
	a bool
	c bool
}`
	if !strings.Contains(string(optimizedContent), expectedOptimizedStructure) {
		t.Errorf("Structure was not optimized as expected. Got:\n%s", string(optimizedContent))
	}
}

// testFilePatterns tests the file pattern matching functionality of gopad.
func testFilePatterns(t *testing.T, dir string) {
	output, err := runCommand("--files", dir, "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("File patterns failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 5") {
		t.Errorf("Unexpected output for file patterns: %s", output)
	}

	output, err = runCommand("--files", dir, "--ignore-pattern", "test\\.go$")
	if err != nil {
		t.Errorf("Ignore patterns failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 4") {
		t.Errorf("Unexpected output for ignore patterns: %s", output)
	}
}

// testErrorHandling tests various error scenarios in gopad.
func testErrorHandling(t *testing.T, dir string) {
	output, err := runCommand()
	if err != nil {
		t.Errorf("Unexpected error when running without arguments: %v", err)
	}
	if !strings.Contains(output, "Usage of gopad:") {
		t.Errorf("Expected usage information, got: %s", output)
	}

	output, err = runCommand("--files", "non_existent_folder")
	if err == nil {
		t.Errorf("Expected error for non-existent folder, got none. Output: %s", output)
	}
	if !strings.Contains(output, "path does not exist") {
		t.Errorf("Expected 'path does not exist' error, got: %s", output)
	}

	output, err = runCommand("--files", dir, "--pattern", "[")
	if err == nil {
		t.Errorf("Expected error for invalid regex, got none. Output: %s", output)
	}
	outputLower := strings.ToLower(output)
	if !strings.Contains(outputLower, "error compiling file pattern regex") {
		t.Errorf("Expected regex compilation error, got: %s", output)
	}
	if !strings.Contains(outputLower, "error parsing regexp: missing closing ]") {
		t.Errorf("Expected specific regex parsing error, got: %s", output)
	}
}

// testUtilityFunctions tests utility functions like version and help in gopad.
func testUtilityFunctions(t *testing.T) {
	output, err := runCommand("--version")
	if err != nil {
		t.Errorf("Version check failed: %v", err)
	}
	if !strings.Contains(output, "Version: ") {
		t.Errorf("Unexpected output for version check: %s", output)
	}

	output, err = runCommand("--help")
	if err != nil {
		t.Errorf("Help check failed: %v", err)
	}
	if !strings.Contains(output, "Usage of gopad:") {
		t.Errorf("Unexpected output for help check: %s", output)
	}
}

// testFlagCombinations tests various combinations of flags in gopad.
func testFlagCombinations(t *testing.T, dir string) {
	output, err := runCommand("--files", dir, "--ignore", filepath.Join(dir, "test.go"), "--view", "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("Flag combinations failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 4") || !strings.Contains(output, filepath.Join(dir, "main.go")) {
		t.Errorf("Unexpected output for flag combinations: %s", output)
	}
}

// testPathsWithSpaces tests gopad's behavior with paths containing spaces.
func testPathsWithSpaces(t *testing.T, dir string) {
	spaceDir := filepath.Join(dir, "folder with spaces")
	err := os.Mkdir(spaceDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory with spaces: %v", err)
	}
	err = os.WriteFile(filepath.Join(spaceDir, "file.go"), []byte("package space"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file in directory with spaces: %v", err)
	}

	output, err := runCommand("--files", spaceDir)
	if err != nil {
		t.Errorf("Paths with spaces failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 1") {
		t.Errorf("Unexpected output for paths with spaces: %s", output)
	}
}

// testRecursiveTraversal tests gopad's recursive directory traversal.
func testRecursiveTraversal(t *testing.T, dir string) {
	output, err := runCommand("--files", dir, "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("Recursive traversal failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 6") {
		t.Errorf("Unexpected output for recursive traversal: %s", output)
	}
}

// testSymbolicLinks tests gopad's handling of symbolic links.
func testSymbolicLinks(t *testing.T, dir string) {
	linkFile := filepath.Join(dir, "link_file.go")
	err := os.Symlink(filepath.Join(dir, "main.go"), linkFile)
	if err != nil {
		t.Fatalf("Failed to create symbolic link: %v", err)
	}

	output, err := runCommand("--files", linkFile)
	if err != nil {
		t.Errorf("Symbolic links failed: %v", err)
	}
	if !strings.Contains(output, "Files analyzed: 1") {
		t.Errorf("Unexpected output for symbolic links: %s", output)
	}
}

// testAccessRights tests gopad's behavior with files having different access rights.
func testAccessRights(t *testing.T, dir string) {
	noAccessFile := filepath.Join(dir, "no_access.go")
	err := os.WriteFile(noAccessFile, []byte("package noaccess"), 0644)
	if err != nil {
		t.Fatalf("Failed to create no-access file: %v", err)
	}

	initialFileCount := countGoFiles(t, dir)

	output, err := runCommand("--files", dir)
	if err != nil {
		t.Errorf("Access rights test failed: %v", err)
	}

	minExpectedFiles := initialFileCount
	maxExpectedFiles := initialFileCount + 1

	actualFileCount := extractFileCount(output)
	if actualFileCount < minExpectedFiles || actualFileCount > maxExpectedFiles {
		t.Errorf("Expected to find between %d and %d files, but got: %d. Full output: %s",
			minExpectedFiles, maxExpectedFiles, actualFileCount, output)
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(noAccessFile, 0000)
		if err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}
		output, _ = runCommand("--files", dir)
		if !strings.Contains(output, "permission denied") {
			t.Errorf("Expected warning about permission denied on non-Windows systems, got: %s", output)
		}
	}
}

// countGoFiles counts the number of Go files in a directory.
func countGoFiles(t *testing.T, dir string) int {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			count++
		}
		return nil
	})
	if err != nil {
		t.Fatalf("Failed to count .go files: %v", err)
	}
	return count
}

// extractFileCount extracts the number of files found from gopad's output.
func extractFileCount(output string) int {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Files analyzed:") {
			var count int
			_, err := fmt.Sscanf(line, "Files analyzed: %d", &count)
			if err == nil {
				return count
			}
		}
	}
	return -1 // Return -1 if unable to extract the number of files
}

func TestMergeFlags(t *testing.T) {
	tests := []struct {
		name     string
		long     string
		short    string
		expected []string
	}{
		{
			name:     "Long flag present",
			long:     "file1.go,file2.go",
			short:    "short.go",
			expected: []string{"file1.go", "file2.go"},
		},
		{
			name:     "Only short flag present",
			long:     "",
			short:    "short1.go,short2.go",
			expected: []string{"short1.go", "short2.go"},
		},
		{
			name:     "Both flags empty",
			long:     "",
			short:    "",
			expected: nil, // Changed from []string{} to nil
		},
		{
			name:     "Long flag with spaces",
			long:     " file1.go , file2.go ",
			short:    "short.go",
			expected: []string{"file1.go", "file2.go"},
		},
		{
			name:     "Short flag with empty parts",
			long:     "",
			short:    "short1.go,,short2.go",
			expected: []string{"short1.go", "short2.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeFlags(tt.long, tt.short)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("mergeFlags(%q, %q) = %v; want %v", tt.long, tt.short, result, tt.expected)
			}
		})
	}
}
func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple split",
			input:    "a,b,c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Split with spaces",
			input:    " a , b , c ",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Split with empty parts",
			input:    "a,,b,c,",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "Single item",
			input:    "a",
			expected: []string{"a"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: nil, // Changed from []string{} to nil
		},
		{
			name:     "Only spaces and commas",
			input:    " , , ",
			expected: nil, // Changed from []string{} to nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitAndTrim(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitAndTrim(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
