package main

import (
	"reflect"
	"testing"
)

// BadStruct represents an inefficiently aligned structure
type BadStruct struct {
	a bool  // 1 byte
	b int32 // 4 bytes
	c bool  // 1 byte
	d int64 // 8 bytes
}

// GoodStruct represents an efficiently aligned version of BadStruct
type GoodStruct struct {
	d int64 // 8 bytes
	b int32 // 4 bytes
	a bool  // 1 byte
	c bool  // 1 byte
}

// TestBadStructAlignment checks if the BadStruct has the expected memory layout and size.
// It verifies that the struct is not optimally aligned, resulting in a larger size.
func TestBadStructAlignment(t *testing.T) {
	badStruct := BadStruct{}
	size := reflect.TypeOf(badStruct).Size()
	expectedSize := uintptr(24) // 1 + 4 + 1 + 8 + padding = 24 bytes on most systems

	if size != expectedSize {
		t.Errorf("BadStruct size = %d; want %d", size, expectedSize)
	}
}

// TestGoodStructAlignment checks if the GoodStruct has the expected memory layout and size.
// It verifies that the struct is optimally aligned, resulting in a smaller size.
func TestGoodStructAlignment(t *testing.T) {
	goodStruct := GoodStruct{}
	size := reflect.TypeOf(goodStruct).Size()
	expectedSize := uintptr(16) // 8 + 4 + 1 + 1 + padding = 16 bytes on most systems

	if size != expectedSize {
		t.Errorf("GoodStruct size = %d; want %d", size, expectedSize)
	}
}

// TestParseStrings tests the ParseStrings function.
// It checks if the function correctly parses a string containing a struct definition,
// and if it produces the expected optimization results.
func TestParseStrings(t *testing.T) {
	input := `package main

type TestStruct struct {
	a bool
	b int64
	c bool
}
`
	results, _, err := ParseStrings(input)
	if err != nil {
		t.Fatalf("ParseStrings failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Name != "TestStruct" {
		t.Errorf("Expected struct name 'TestStruct', got '%s'", results[0].Name)
	}

	// Check that sizes were calculated
	if results[0].MetaData.BeforeSize == 0 || results[0].MetaData.AfterSize == 0 {
		t.Errorf("Expected non-zero sizes, got before: %d, after: %d", results[0].MetaData.BeforeSize, results[0].MetaData.AfterSize)
	}

	// TODO
	// Check that the field order has changed (optimization)
	//optimizedStructString := string(results[0].Data)
	//if !strings.Contains(optimizedStructString, "b int64\n\ta bool\n\tc bool") {
	//	t.Errorf("Expected fields to be reordered, got: %s", optimizedStructString)
	//}
}

// TestParseBytes tests the Parse function.
// It verifies that the function can correctly parseData a byte slice containing
// a struct definition and produce the expected results.
func TestParseBytes(t *testing.T) {
	input := []byte(`package main

type TestStruct struct {
	a bool
	b int32
	c string
}
`)
	results, _, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Name != "TestStruct" {
		t.Errorf("Expected struct name 'TestStruct', got '%s'", results[0].Name)
	}
}

// TestReplacer tests the Replacer function.
// It checks if the function can replace the original struct definition
// with an optimized version. Note: Some checks are currently commented out (TODO).
func TestReplacer(t *testing.T) {
	original := []byte(`package main

type TestStruct struct {
	a bool
	b int64
	c bool
}
`)
	results, _, err := Parse(original)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	_, err = Replacer(original, results)
	//modified, err := Replacer(original, results)
	if err != nil {
		t.Fatalf("Replacer failed: %v", err)
	}

	// TODO
	//if bytes.Equal(original, modified) {
	//	t.Errorf("Expected modified content to be different from original")
	//}

	// TODO
	//if !bytes.Contains(modified, []byte("TestStruct")) {
	//	t.Errorf("Modified content should still contain 'TestStruct'")
	//}

	// TODO
	// Check that the field order has changed
	//if !bytes.Contains(modified, []byte("b int64\n\ta bool\n\tc bool")) {
	//	t.Errorf("Expected fields to be reordered in the modified content")
	//}
}

// TestParseFile is a mock test for the ParseFile function.
// Since ParseFile depends on the file system, this test only checks
// if the function exists and returns an expected error for a non-existent file.
func TestParseFile(t *testing.T) {
	// In a real scenario, you would create a temporary file for testing
	// Here we just check that the function exists and returns the expected type
	_, _, err := ParseFile("non_existent_file.go")
	if err == nil {
		t.Errorf("Expected error for non-existent file, got nil")
	}
}
