package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	textreplacer "github.com/t34-dev/go-text-replacer"
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
	// Normalize line endings to LF
	bytes = normalizeLineEndings(bytes)

	node, err := parser.ParseFile(token.NewFileSet(), path, bytes, parser.ParseComments)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Failed to parseData source: %v", err))
	}

	//var results []MetaData
	var structures []*Structure
	mapperItems := map[string]*Structure{}

	ast.Inspect(node, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		// Don't subtract len("type ") here - this leads to incorrect start pos
		// in cases when struct is in a "type" block (type (...)).
		startPos := int(typeSpec.Pos())
		plus := 0
		if typeSpec.Comment != nil && len(typeSpec.Comment.List) > 0 {
			plus += len(typeSpec.Comment.List[0].Text) + 1
		}
		endPos := int(typeSpec.Type.End()) + plus
		metaData := MetaData{
			StartPos: startPos,
			EndPos:   endPos,
		}
		item := createTypeItemInfo(typeSpec, nil, mapperItems)
		item.MetaData = &metaData
		if item != nil {
			structures = append(structures, item)
		}
		return true
	})
	return structures, mapperItems, err
}

func createMapperItem(structure *Structure, mapperItems map[string]*Structure) map[string]*Structure {
	mapperItems[structure.Path] = structure
	if structure.IsStructure {
		for _, elem := range structure.NestedFields {
			createMapperItem(elem, mapperItems)
		}
	}
	return mapperItems
}

func createMapper(structures []*Structure) map[string]*Structure {
	mapper := map[string]*Structure{}
	for _, structure := range structures {
		createMapperItem(structure, mapper)
	}
	return mapper
}

// Replacer replaces the original struct definitions with optimized versions in the source code
func Replacer(file []byte, structures []*Structure) ([]byte, error) {
	// Normalize line endings to LF
	file = normalizeLineEndings(file)

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
