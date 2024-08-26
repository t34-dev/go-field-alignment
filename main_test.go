package main

import (
	"bytes"
	"os"
	"testing"
)

// testEnterFile is the path to the input test file
const testEnterFile = "tests/enter/file.go"

// testOutFile is the path to the expected output test file
const testOutFile = "tests/out/file.go"

// TestStructAlignment tests the alignment and optimization of struct fields.
// It reads input and expected output files, applies the optimization,
// and compares the result with the expected output.
func TestStructAlignment(t *testing.T) {
	// read files
	enterFIle, err := os.ReadFile(testEnterFile)
	if err != nil {
		t.Fatal(err)
	}
	outFIle, err := os.ReadFile(testOutFile)
	if err != nil {
		t.Fatal(err)
	}

	// parser
	structures, mapper, err := Parse(enterFIle)
	if err != nil {
		t.Fatal(err)
	}
	calculateStructures(structures, true)
	debugPrintStructures(structures)

	optimizeMapperStructures(mapper)
	calculateStructures(structures, false)
	debugPrintStructures(structures)

	err = renderStructures(structures)
	if err != nil {
		t.Fatal(err)
	}
	// Replace content
	resultFile, err := Replacer(enterFIle, structures)
	if err != nil {
		t.Fatal(err)
	}

	// Compare modified code with structsSourceOut
	if !bytes.Equal(resultFile, outFIle) {
		t.Errorf("Modified code does not match expected output.\nGot:\n%s\nWant:\n%s", string(resultFile), string(outFIle))
	} else {
		t.Log("Modified code matches expected output.")
	}
}
