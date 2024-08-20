package main

import (
	"go/ast"
)

type BadStruct struct {
	a bool  // 1 byte
	b int32 // 4 bytes
	c bool  // 1 byte
	d int64 // 8 bytes
}

type GoodStruct struct {
	d int64 // 8 bytes
	b int32 // 4 bytes
	a bool  // 1 byte
	c bool  // 1 byte
}

type ItemInfo struct {
	Name         string
	Path         string
	Root         *ast.TypeSpec
	RootField    *ast.Field
	StructType   ast.Expr
	StringType   string
	IsStructure  bool
	Size         uintptr
	Align        uintptr
	Offset       uintptr
	NestedFields []*ItemInfo
}
