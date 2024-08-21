package main

import (
	"errors"
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"go/ast"
	"go/parser"
	"go/token"
)

// ViewMode is a global flag to control the output verbosity
var (
	ViewMode = false
)

// ItemInfo represents detailed information about a struct field or type
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

// Result represents the outcome of struct optimization
type Result struct {
	NameStructure string
	BeforeSize    uintptr
	AfterSize     uintptr
	Data          []byte
	StartPos      int
	EndPos        int
}

// ParseFile parses a Go file and returns optimization results
func ParseFile(path string) ([]Result, error) {
	return parse(path, nil)
}

// ParseBytes parses Go code from a byte slice and returns optimization results
func ParseBytes(bytes []byte) ([]Result, error) {
	return parse("", bytes)
}

// ParseStrings parses Go code from a string and returns optimization results
func ParseStrings(str string) ([]Result, error) {
	return parse("", []byte(str))
}

// parse is the core function that handles parsing and optimization of Go code
func parse(path string, bytes []byte) ([]Result, error) {
	node, err := parser.ParseFile(token.NewFileSet(), path, bytes, parser.ParseComments)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to parse source: %v", err))
	}

	var results []Result
	var structures []*ItemInfo
	mapperItems := map[string]*ItemInfo{}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			startPos := int(typeSpec.Pos()) - len("type ")
			plus := 0
			if typeSpec.Comment != nil && len(typeSpec.Comment.List) > 0 {
				plus += len(typeSpec.Comment.List[0].Text) + 1
			}
			endPos := int(typeSpec.Type.End()) + plus
			results = append(results, Result{
				StartPos: startPos,
				EndPos:   endPos,
			})
			item := createItemInfo(typeSpec, nil, mapperItems)
			if item != nil {
				structures = append(structures, item)
			}
		}
		return true
	})
	calculateStructures(structures)
	for idx, structure := range structures {
		results[idx].BeforeSize = structure.Size
	}
	if ViewMode {
		viewPrintStructures(structures)
	}

	// Optimize structures
	optimizeStructures(mapperItems)
	calculateStructures(structures)
	for idx, structure := range structures {
		results[idx].AfterSize = structure.Size
	}
	if ViewMode {
		viewPrintStructures(structures)
	}

	// add data
	for idx, structure := range structures {
		code, err := formatGoCode(renderStructure(structure))
		if err != nil {
			return nil, err
		}
		results[idx].NameStructure = structure.Name
		results[idx].Data = []byte(code)
	}
	return results, nil
}

// Replacer replaces the original struct definitions with optimized versions in the source code
func Replacer(file []byte, results []Result) ([]byte, error) {
	var blocks []textreplacer.Block
	for _, elem := range results {
		blocks = append(blocks, textreplacer.Block{
			Start: elem.StartPos - 1,
			End:   elem.EndPos - 1,
			Txt:   elem.Data,
		})
	}
	replacer := textreplacer.New(file)
	return replacer.Enter(blocks)
}
