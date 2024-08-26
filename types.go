package main

import (
	"errors"
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"go/ast"
	"go/parser"
	"go/token"
)

// MetaData represents the outcome of struct optimization
type MetaData struct {
	BeforeSize uintptr
	AfterSize  uintptr
	Data       []byte
	StartPos   int
	EndPos     int
}

// Structure represents detailed information about a struct field or type
type Structure struct {
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
	NestedFields []*Structure
	MetaData     *MetaData
}

// ParseFile parses a Go file and returns optimization results
func ParseFile(path string) ([]*Structure, map[string]*Structure, error) {
	return parseData(path, nil)
}

// Parse parses Go code from a byte slice and returns optimization results
func Parse(bytes []byte) ([]*Structure, map[string]*Structure, error) {
	return parseData("", bytes)
}

// ParseStrings parses Go code from a string and returns optimization results
func ParseStrings(str string) ([]*Structure, map[string]*Structure, error) {
	return parseData("", []byte(str))
}

// parseData is the core function that handles parsing and optimization of Go code
func parseData(path string, bytes []byte) ([]*Structure, map[string]*Structure, error) {
	node, err := parser.ParseFile(token.NewFileSet(), path, bytes, parser.ParseComments)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to parseData source: %v", err))
	}

	//var results []MetaData
	var structures []*Structure
	mapperItems := map[string]*Structure{}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			startPos := int(typeSpec.Pos()) - len("type ")
			plus := 0
			if typeSpec.Comment != nil && len(typeSpec.Comment.List) > 0 {
				plus += len(typeSpec.Comment.List[0].Text) + 1
			}
			endPos := int(typeSpec.Type.End()) + plus
			metaData := MetaData{
				StartPos: startPos,
				EndPos:   endPos,
			}
			item := createItemInfo(typeSpec, nil, mapperItems)
			item.MetaData = &metaData
			if item != nil {
				structures = append(structures, item)
			}
		}
		return true
	})
	return structures, mapperItems, err
}

// Replacer replaces the original struct definitions with optimized versions in the source code
func Replacer(file []byte, structures []*Structure) ([]byte, error) {
	var blocks []textreplacer.Block
	for _, elem := range structures {
		blocks = append(blocks, textreplacer.Block{
			Start: elem.MetaData.StartPos - 1,
			End:   elem.MetaData.EndPos - 1,
			Txt:   elem.MetaData.Data,
		})
	}
	replacer := textreplacer.New(file)
	return replacer.Enter(blocks)
}
