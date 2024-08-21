package main

import (
	"errors"
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"go/ast"
	"go/parser"
	"go/token"
)

var (
	ViewMode = false
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

type Result struct {
	NameStructure string
	beforeSize    int
	AfterSize     int
	Data          []byte
	StartPos      int
	EndPos        int
}

func ParseFile(path string) ([]Result, error) {
	return parse(path, nil)
}
func ParseBytes(bytes []byte) ([]Result, error) {
	return parse("", bytes)
}
func ParseStrings(str string) ([]Result, error) {
	return parse("", []byte(str))
}
func parse(path string, bytes []byte) ([]Result, error) {
	node, err := parser.ParseFile(token.NewFileSet(), path, bytes, parser.ParseComments)
	if err != nil {
		errors.New(fmt.Sprintf("Failed to parse source: %v", err))
	}

	var structures []*ItemInfo
	mapperItems := map[string]*ItemInfo{}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			item := createItemInfo(typeSpec, nil, mapperItems)
			if item != nil {
				structures = append(structures, item)
			}
		}
		return true
	})
	calculateStructures(structures)
	if ViewMode {
		viewPrintStructures(structures)
	}

	// Optimize structures
	optimizeStructures(mapperItems)
	calculateStructures(structures)
	if ViewMode {
		viewPrintStructures(structures)
	}

	var results []Result
	for _, structure := range structures {
		code, err := formatGoCode(renderStructure(structure))
		if err != nil {
			return nil, err
		}
		results = append(results, Result{
			beforeSize:    1,
			AfterSize:     2,
			NameStructure: structure.Name,
			Data:          []byte(code),
		})
	}
	return results, nil
}
func Replacer(file []byte, results []Result) ([]byte, error) {
	var blocks []textreplacer.Block
	for _, elem := range results {
		blocks = append(blocks, textreplacer.Block{
			Start: elem.StartPos,
			End:   elem.EndPos,
			Txt:   elem.Data,
		})
	}
	replacer := textreplacer.New(file)
	return replacer.Enter(blocks)
}
