package main

import (
	"go/ast"
	"reflect"
	"testing"
)

// TestFormatGoCode tests the formatGoCode function.
// It checks if the function correctly formats valid Go code
// and handles malformed code appropriately.
func TestFormatGoCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple function",
			input:    "func main() { fmt.Println(\"Hello, World!\") }",
			expected: "func main() { fmt.Println(\"Hello, World!\") }",
		},
		{
			name:     "malformed code",
			input:    "func main() { fmt.Println(\"Hello, World!\") ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := formatGoCode(tt.input)
			if tt.expected == "" {
				if err == nil {
					t.Errorf("Expected error for malformed code, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}

// TestGetTypeString tests the getTypeString function.
// It verifies that the function correctly converts various AST expressions
// to their string representations.
func TestGetTypeString(t *testing.T) {
	tests := []struct {
		name     string
		expr     ast.Expr
		expected string
	}{
		{
			name:     "identifier",
			expr:     &ast.Ident{Name: "int"},
			expected: "int",
		},
		{
			name:     "pointer",
			expr:     &ast.StarExpr{X: &ast.Ident{Name: "int"}},
			expected: "*int",
		},
		{
			name:     "slice",
			expr:     &ast.ArrayType{Elt: &ast.Ident{Name: "int"}},
			expected: "[]int",
		},
		{
			name: "map",
			expr: &ast.MapType{
				Key:   &ast.Ident{Name: "string"},
				Value: &ast.Ident{Name: "int"},
			},
			expected: "map[string]int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTypeString(tt.expr)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestGetFuncTypeString tests the getFuncTypeString function.
// It checks if the function correctly generates string representations
// of function types with various parameters and return values.
func TestGetFuncTypeString(t *testing.T) {
	tests := []struct {
		name     string
		funcType *ast.FuncType
		expected string
	}{
		{
			name: "no params, no results",
			funcType: &ast.FuncType{
				Params:  &ast.FieldList{},
				Results: nil,
			},
			expected: "func()",
		},
		{
			name: "with params and results",
			funcType: &ast.FuncType{
				Params: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "int"}},
						{Type: &ast.Ident{Name: "string"}},
					},
				},
				Results: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "bool"}},
					},
				},
			},
			expected: "func(int, string) bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFuncTypeString(tt.funcType)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSortMapKeysBySlashCount tests the sortMapKeysBySlashCount function.
// It verifies that the function correctly sorts ItemInfo slices
// based on the number of slashes in their paths.
func TestSortMapKeysBySlashCount(t *testing.T) {
	input := map[string]*ItemInfo{
		"a":     {Path: "a"},
		"a/b":   {Path: "a/b"},
		"a/b/c": {Path: "a/b/c"},
		"x/y":   {Path: "x/y"},
	}

	expected := []*ItemInfo{
		{Path: "a/b/c"},
		{Path: "a/b"},
		{Path: "x/y"},
		{Path: "a"},
	}

	result := sortMapKeysBySlashCount(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// TestMaxValue tests the maxValue function.
// It checks if the function correctly returns the maximum of two integers
// for various input combinations.
func TestMaxValue(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{-1, 0, 0},
		{10, 10, 10},
	}

	for _, tt := range tests {
		result := maxValue(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("maxValue(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

// TestIsValidCustomTypeName tests the isValidCustomTypeName function.
// It verifies that the function correctly identifies valid and invalid
// custom type names according to Go naming conventions.
func TestIsValidCustomTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		{"valid custom type", "MyType", true},
		{"valid with underscore", "My_Type", true},
		{"valid with number", "Type2", true},
		{"invalid starts with number", "2Type", false},
		{"invalid contains space", "My Type", false},
		{"invalid standard type", "int", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidCustomTypeName(tt.typeName)
			if result != tt.expected {
				t.Errorf("isValidCustomTypeName(%q) = %v; want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}
