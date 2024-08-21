package main

import (
	"fmt"
	"go/ast"
	"unsafe"
)

// getTypeID generates a unique identifier for an AST expression.
// This is used to detect recursive types and prevent infinite loops.
func getTypeID(expr ast.Expr) string {
	return fmt.Sprintf("%T:%p", expr, expr)
}

// getFieldSizeWithMap calculates the size of a field in a structure.
// It handles various types including basic types, pointers, arrays, structs, maps, channels, and interfaces.
// The function uses a map to keep track of seen types to handle recursive structures.
func getFieldSizeWithMap(field ast.Expr, seenTypes map[string]bool) uintptr {
	typeID := getTypeID(field)

	if seenTypes[typeID] {
		return unsafe.Sizeof(uintptr(0))
	}

	seenTypes[typeID] = true
	defer delete(seenTypes, typeID)

	switch t := (field).(type) {
	case *ast.Ident:
		switch t.Name {
		case "bool":
			return unsafe.Sizeof(bool(false))
		case "int8", "uint8", "byte":
			return unsafe.Sizeof(int8(0))
		case "int16", "uint16":
			return unsafe.Sizeof(int16(0))
		case "int32", "uint32", "float32", "rune":
			return unsafe.Sizeof(int32(0))
		case "int64", "uint64", "float64":
			return unsafe.Sizeof(int64(0))
		case "int", "uint":
			return unsafe.Sizeof(int(0))
		case "string":
			return unsafe.Sizeof("")
		case "complex64":
			return unsafe.Sizeof(complex64(0))
		case "complex128":
			return unsafe.Sizeof(complex128(0))
		}
	case *ast.StarExpr:
		return unsafe.Sizeof(uintptr(0))
	case *ast.ArrayType:
		if t.Len == nil {
			return unsafe.Sizeof([]int{})
		} else {
			elemSize := getFieldSizeWithMap(t.Elt, seenTypes)
			length := 0
			if lit, ok := t.Len.(*ast.BasicLit); ok {
				fmt.Sscanf(lit.Value, "%d", &length)
			}
			return elemSize * uintptr(length) // Remove padding
		}
	case *ast.StructType:
		var size, maxAlign uintptr
		for _, field := range t.Fields.List {
			fieldSize := getFieldSizeWithMap(field.Type, seenTypes)
			fieldAlign := getFieldAlign(field.Type)
			size = align(size, fieldAlign) + fieldSize
			if fieldAlign > maxAlign {
				maxAlign = fieldAlign
			}
		}
		return align(size, maxAlign)
	case *ast.MapType:
		return unsafe.Sizeof(map[string]int{})
	case *ast.ChanType:
		return unsafe.Sizeof(make(chan int))
	case *ast.InterfaceType:
		return unsafe.Sizeof((*interface{})(nil))
	}
	return unsafe.Sizeof("")
}

// getFieldAlign determines the alignment requirement of a field.
// It handles various types similar to getFieldSizeWithMap.
func getFieldAlign(field ast.Expr) uintptr {
	switch t := (field).(type) {
	case *ast.Ident:
		switch t.Name {
		case "bool":
			return unsafe.Alignof(false)
		case "int8", "uint8", "byte":
			return unsafe.Alignof(int8(0))
		case "int16", "uint16":
			return unsafe.Alignof(int16(0))
		case "int32", "uint32", "float32", "rune":
			return unsafe.Alignof(int32(0))
		case "int64", "uint64", "float64":
			return unsafe.Alignof(int64(0))
		case "int", "uint":
			return unsafe.Alignof(0)
		case "string":
			return unsafe.Alignof("")
		}
	case *ast.StarExpr:
		return unsafe.Alignof(uintptr(0))
	case *ast.ArrayType:
		return getFieldAlign(t.Elt)
	case *ast.StructType:
		var maxAlign uintptr
		for _, field := range t.Fields.List {
			fieldAlign := getFieldAlign(field.Type)
			if fieldAlign > maxAlign {
				maxAlign = fieldAlign
			}
		}
		return maxAlign
	case *ast.MapType:
		return unsafe.Alignof(map[string]int{})
	case *ast.ChanType:
		return unsafe.Alignof(make(chan int))
	case *ast.InterfaceType:
		return unsafe.Alignof((*interface{})(nil))
	}
	return unsafe.Alignof("")
}

// align calculates the next aligned address given a size and an alignment.
// This function is used to ensure proper alignment of fields within a structure.
func align(size, align uintptr) uintptr {
	return (size + align - 1) &^ (align - 1)
}

// getFieldSize is a wrapper function that initializes a new map and calls getFieldSizeWithMap.
// This function is the main entry point for calculating field sizes.
func getFieldSize(field ast.Expr) uintptr {
	return getFieldSizeWithMap(field, make(map[string]bool))
}
