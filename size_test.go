package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"
	"unsafe"
)

// TestGetFieldSizeStruct tests the getFieldSize function with a simple struct.
// It compares the calculated size with the actual size of an equivalent Go struct.
func TestGetFieldSizeStruct(t *testing.T) {
	structExpr := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{Type: &ast.Ident{Name: "int"}},
				{Type: &ast.Ident{Name: "string"}},
				{Type: &ast.Ident{Name: "bool"}},
			},
		},
	}

	var expr ast.Expr = structExpr
	size := getFieldSize(expr)

	// Create a real structure for comparison
	type testStruct struct {
		i int
		s string
		b bool
	}
	expected := reflect.TypeOf(testStruct{}).Size()

	// Output detailed information about sizes and alignment
	t.Logf("Size of int: %d", unsafe.Sizeof(int(0)))
	t.Logf("Size of string: %d", unsafe.Sizeof(""))
	t.Logf("Size of bool: %d", unsafe.Sizeof(false))
	t.Logf("Structure alignment: %d", reflect.TypeOf(testStruct{}).Align())

	v := reflect.ValueOf(testStruct{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		t.Logf("Field %s: size %d, offset %d, alignment %d",
			field.Name, field.Type.Size(), field.Offset, field.Type.Align())
	}

	t.Logf("Expected size: %d", expected)
	t.Logf("Actual size: %d", size)

	if size != expected {
		t.Errorf("getFieldSize() for struct = %v, want %v", size, expected)
	}
}

// TestGetFieldSizeTypes tests the getFieldSize function for various basic types.
// It checks if the calculated sizes match the actual sizes of Go types.
func TestGetFieldSizeTypes(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{"bool", &ast.Ident{Name: "bool"}, unsafe.Sizeof(bool(false))},
		{"int", &ast.Ident{Name: "int"}, unsafe.Sizeof(int(0))},
		{"int8", &ast.Ident{Name: "int8"}, unsafe.Sizeof(int8(0))},
		{"int16", &ast.Ident{Name: "int16"}, unsafe.Sizeof(int16(0))},
		{"int32", &ast.Ident{Name: "int32"}, unsafe.Sizeof(int32(0))},
		{"int64", &ast.Ident{Name: "int64"}, unsafe.Sizeof(int64(0))},
		{"uint", &ast.Ident{Name: "uint"}, unsafe.Sizeof(uint(0))},
		{"uint8", &ast.Ident{Name: "uint8"}, unsafe.Sizeof(uint8(0))},
		{"uint16", &ast.Ident{Name: "uint16"}, unsafe.Sizeof(uint16(0))},
		{"uint32", &ast.Ident{Name: "uint32"}, unsafe.Sizeof(uint32(0))},
		{"uint64", &ast.Ident{Name: "uint64"}, unsafe.Sizeof(uint64(0))},
		{"float32", &ast.Ident{Name: "float32"}, unsafe.Sizeof(float32(0))},
		{"float64", &ast.Ident{Name: "float64"}, unsafe.Sizeof(float64(0))},
		{"complex64", &ast.Ident{Name: "complex64"}, unsafe.Sizeof(complex64(0))},
		{"complex128", &ast.Ident{Name: "complex128"}, unsafe.Sizeof(complex128(0))},
		{"string", &ast.Ident{Name: "string"}, unsafe.Sizeof("")},
		{"slice", &ast.ArrayType{Elt: &ast.Ident{Name: "int"}}, unsafe.Sizeof([]int{})},
		{"array", &ast.ArrayType{Elt: &ast.Ident{Name: "int"}, Len: &ast.BasicLit{Kind: token.INT, Value: "5"}}, unsafe.Sizeof([5]int{})},
		{"map", &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}}, unsafe.Sizeof(map[string]int{})},
		{"chan", &ast.ChanType{Value: &ast.Ident{Name: "int"}}, unsafe.Sizeof(make(chan int))},
		{"interface", &ast.InterfaceType{}, unsafe.Sizeof((*interface{})(nil))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFieldSize(tt.expr)
			if got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFieldSizePointers tests the getFieldSize function for pointer types.
// It verifies that the function correctly calculates sizes for single and double pointers.
func TestGetFieldSizePointers(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{"*int", &ast.StarExpr{X: &ast.Ident{Name: "int"}}, unsafe.Sizeof((*int)(nil))},
		{"*string", &ast.StarExpr{X: &ast.Ident{Name: "string"}}, unsafe.Sizeof((*string)(nil))},
		{"*bool", &ast.StarExpr{X: &ast.Ident{Name: "bool"}}, unsafe.Sizeof((*bool)(nil))},
		{"**int", &ast.StarExpr{X: &ast.StarExpr{X: &ast.Ident{Name: "int"}}}, unsafe.Sizeof((**int)(nil))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldSize(tt.expr); got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFieldSizeComplexStructs tests the getFieldSize function for more complex struct types.
// It includes tests for structs with embedded structs, slices, and maps.
func TestGetFieldSizeComplexStructs(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{
			name: "struct with embedded struct",
			expr: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "int"}},
						{Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "string"}},
									{Type: &ast.Ident{Name: "bool"}},
								},
							},
						}},
					},
				},
			},
			want: reflect.TypeOf(struct {
				i int
				s struct {
					str string
					b   bool
				}
			}{}).Size(),
		},
		{
			name: "struct with slice and map",
			expr: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}}},
						{Type: &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}}},
					},
				},
			},
			want: reflect.TypeOf(struct {
				s []int
				m map[string]int
			}{}).Size(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldSize(tt.expr); got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetFieldSizeRecursiveStruct tests the getFieldSize function with a recursive struct.
// It verifies that the function can handle recursive types without infinite loops.
func TestGetFieldSizeRecursiveStruct(t *testing.T) {
	// Define a recursive structure
	type RecursiveStruct struct {
		i int
		r *RecursiveStruct
	}

	// Create AST representation of RecursiveStruct
	recursiveStructExpr := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: "i"}},
					Type:  &ast.Ident{Name: "int"},
				},
				{
					Names: []*ast.Ident{{Name: "r"}},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "RecursiveStruct"},
					},
				},
			},
		},
	}

	// Create full AST representation of the type definition
	typeSpec := &ast.TypeSpec{
		Name: &ast.Ident{Name: "RecursiveStruct"},
		Type: recursiveStructExpr,
	}

	// Wrap the type definition in GenDecl, as it would be in a real Go file
	genDecl := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{typeSpec},
	}

	// Now we have a complete AST representation of the RecursiveStruct type definition

	var expr ast.Expr = recursiveStructExpr
	size := getFieldSize(expr)

	expected := reflect.TypeOf(RecursiveStruct{}).Size()

	t.Logf("AST representation of the structure:")
	t.Logf("%#v", genDecl)
	t.Logf("Size of int: %d", unsafe.Sizeof(int(0)))
	t.Logf("Size of pointer: %d", unsafe.Sizeof(uintptr(0)))
	t.Logf("Expected size: %d", expected)
	t.Logf("Actual size: %d", size)

	if size != expected {
		t.Errorf("getFieldSize() for recursive struct = %v, want %v", size, expected)
	}
}
