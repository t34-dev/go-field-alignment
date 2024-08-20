package main

import (
	"fmt"
	"go/ast"
	"unsafe"
)

func getTypeID(expr ast.Expr) string {
	return fmt.Sprintf("%T:%p", expr, expr)
}
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
			return elemSize * uintptr(length) // Delete padding
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

func align(size, align uintptr) uintptr {
	return (size + align - 1) &^ (align - 1)
}

func getFieldSize(field ast.Expr) uintptr {
	return getFieldSizeWithMap(field, make(map[string]bool))
}
