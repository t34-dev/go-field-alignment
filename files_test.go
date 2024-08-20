package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGoPadding(t *testing.T) {
	// Setup: Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "go-padding-test")
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

func runCommand(args ...string) (string, error) {
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func testBasicUsage(t *testing.T, dir string) {
	output, err := runCommand("--file", filepath.Join(dir, "main.go"))
	if err != nil {
		t.Errorf("Basic usage failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 1") {
		t.Errorf("Unexpected output for basic usage: %s", output)
	}

	output, err = runCommand("-f", fmt.Sprintf("%s,%s", filepath.Join(dir, "main.go"), filepath.Join(dir, "utils.go")))
	if err != nil {
		t.Errorf("Basic usage with short flag failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 2") {
		t.Errorf("Unexpected output for basic usage with short flag: %s", output)
	}
}

func testIgnoreFiles(t *testing.T, dir string) {
	output, err := runCommand("--file", dir, "--ignore", filepath.Join(dir, "test.go")+","+filepath.Join(dir, "ignore_me.go"))
	if err != nil {
		t.Errorf("Ignore files failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 3") {
		t.Errorf("Unexpected output for ignore files: %s", output)
	}

	output, err = runCommand("-f", dir, "-i", fmt.Sprintf("%s,%s", filepath.Join(dir, "test.go"), filepath.Join(dir, "ignore_me.go")))
	if err != nil {
		t.Errorf("Ignore files with short flag failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 3") {
		t.Errorf("Unexpected output for ignore files with short flag: %s", output)
	}
}

func testViewFiles(t *testing.T, dir string) {
	output, err := runCommand("--file", dir, "--view")
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

func testApplyFixes(t *testing.T, dir string) {
	output, err := runCommand("--file", filepath.Join(dir, "main.go"), "--fix")
	if err != nil {
		t.Errorf("Apply fixes failed: %v", err)
	}
	if !strings.Contains(output, "Applying fixes to files:") {
		t.Errorf("Unexpected output for apply fixes: %s", output)
	}
}

func testFilePatterns(t *testing.T, dir string) {
	output, err := runCommand("--file", dir, "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("File patterns failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 5") {
		t.Errorf("Unexpected output for file patterns: %s", output)
	}

	output, err = runCommand("--file", dir, "--ignore-pattern", "test\\.go$")
	if err != nil {
		t.Errorf("Ignore patterns failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 4") {
		t.Errorf("Unexpected output for ignore patterns: %s", output)
	}
}

func testErrorHandling(t *testing.T, dir string) {
	output, err := runCommand()
	if err != nil {
		t.Errorf("Unexpected error when running without arguments: %v", err)
	}
	if !strings.Contains(output, "Usage of go-padding:") {
		t.Errorf("Expected usage information, got: %s", output)
	}

	output, err = runCommand("--file", "non_existent_folder")
	if err == nil {
		t.Errorf("Expected error for non-existent folder, got none. Output: %s", output)
	}
	if !strings.Contains(output, "path does not exist") {
		t.Errorf("Expected 'path does not exist' error, got: %s", output)
	}

	output, err = runCommand("--file", dir, "--pattern", "[")
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
	if !strings.Contains(output, "Usage of go-padding:") {
		t.Errorf("Unexpected output for help check: %s", output)
	}
}

func testFlagCombinations(t *testing.T, dir string) {
	output, err := runCommand("--file", dir, "--ignore", filepath.Join(dir, "test.go"), "--view", "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("Flag combinations failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 4") || !strings.Contains(output, filepath.Join(dir, "main.go")) {
		t.Errorf("Unexpected output for flag combinations: %s", output)
	}
}

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

	output, err := runCommand("--file", spaceDir)
	if err != nil {
		t.Errorf("Paths with spaces failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 1") {
		t.Errorf("Unexpected output for paths with spaces: %s", output)
	}
}

func testRecursiveTraversal(t *testing.T, dir string) {
	output, err := runCommand("--file", dir, "--pattern", "\\.go$")
	if err != nil {
		t.Errorf("Recursive traversal failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 6") {
		t.Errorf("Unexpected output for recursive traversal: %s", output)
	}
}

func testSymbolicLinks(t *testing.T, dir string) {
	linkFile := filepath.Join(dir, "link_file.go")
	err := os.Symlink(filepath.Join(dir, "main.go"), linkFile)
	if err != nil {
		t.Fatalf("Failed to create symbolic link: %v", err)
	}

	output, err := runCommand("--file", linkFile)
	if err != nil {
		t.Errorf("Symbolic links failed: %v", err)
	}
	if !strings.Contains(output, "Files found: 1") {
		t.Errorf("Unexpected output for symbolic links: %s", output)
	}
}

func testAccessRights(t *testing.T, dir string) {
	noAccessFile := filepath.Join(dir, "no_access.go")
	err := os.WriteFile(noAccessFile, []byte("package noaccess"), 0644)
	if err != nil {
		t.Fatalf("Failed to create no-access file: %v", err)
	}

	initialFileCount := countGoFiles(t, dir)

	if runtime.GOOS != "windows" {
		err = os.Chmod(noAccessFile, 0000)
		if err != nil {
			t.Fatalf("Failed to change file permissions: %v", err)
		}
	}

	output, err := runCommand("--file", dir)
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
		if !strings.Contains(output, "Warning: Permission denied") {
			t.Errorf("Expected warning about permission denied on non-Windows systems, got: %s", output)
		}
	}
}

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

func extractFileCount(output string) int {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "Files found:") {
			var count int
			_, err := fmt.Sscanf(line, "Files found: %d", &count)
			if err == nil {
				return count
			}
		}
	}
	return -1 // Возвращаем -1, если не удалось извлечь количество файлов
}
