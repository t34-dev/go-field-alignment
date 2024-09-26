package main

import (
	"fmt"
	"go/ast"
	"sort"
	"strings"
	"unicode"
)

// stdTypes is a map of standard Go types used for type checking
var stdTypes = map[string]bool{
	"bool":       true,
	"string":     true,
	"int":        true,
	"int8":       true,
	"int16":      true,
	"int32":      true,
	"int64":      true,
	"uint":       true,
	"uint8":      true,
	"uint16":     true,
	"uint32":     true,
	"uint64":     true,
	"uintptr":    true,
	"byte":       true,
	"rune":       true,
	"float32":    true,
	"float64":    true,
	"complex64":  true,
	"complex128": true,
}

// getTypeString returns a string representation of an AST expression
func getTypeString(expr ast.Expr) string {
	if expr == nil {
		return ""
	}

	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.BasicLit:
		return t.Value
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			return "[]" + getTypeString(t.Elt)
		}
		return fmt.Sprintf("[%s]%s", getTypeString(t.Len), getTypeString(t.Elt))
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	case *ast.FuncType:
		return getFuncTypeString(t)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", getTypeString(t.Key), getTypeString(t.Value))
	case *ast.ChanType:
		return getChanTypeString(t)
	case *ast.StructType:
		return "struct{...}"
	case *ast.InterfaceType:
		return "interface{...}"
	case *ast.Ellipsis:
		return "..." + getTypeString(t.Elt)
	case *ast.ParenExpr:
		return "(" + getTypeString(t.X) + ")"
	case *ast.CompositeLit:
		return getTypeString(t.Type)
	case *ast.FuncLit:
		return getFuncTypeString(t.Type)
	default:
		return fmt.Sprintf("%T", expr)
	}
}

// getFuncTypeString returns a string representation of a function type
func getFuncTypeString(t *ast.FuncType) string {
	params := getFieldListString(t.Params)
	results := getFieldListString(t.Results)

	if results == "" {
		return fmt.Sprintf("func(%s)", params)
	}
	return fmt.Sprintf("func(%s) %s", params, results)
}

// getChanTypeString returns a string representation of a channel type
func getChanTypeString(t *ast.ChanType) string {
	switch t.Dir {
	case ast.SEND:
		return fmt.Sprintf("chan<- %s", getTypeString(t.Value))
	case ast.RECV:
		return fmt.Sprintf("<-chan %s", getTypeString(t.Value))
	default:
		return fmt.Sprintf("chan %s", getTypeString(t.Value))
	}
}

// getFieldListString returns a string representation of a field list
func getFieldListString(fields *ast.FieldList) string {
	if fields == nil {
		return ""
	}
	var parts []string
	for _, field := range fields.List {
		typeStr := getTypeString(field.Type)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				parts = append(parts, fmt.Sprintf("%s %s", name.Name, typeStr))
			}
		} else {
			parts = append(parts, typeStr)
		}
	}
	return strings.Join(parts, ", ")
}

// sortMapKeysBySlashCount sorts Structure slices by their path's slash count
func sortMapKeysBySlashCount(inputMap map[string]*Structure) []*Structure {
	items := make([]*Structure, 0, len(inputMap))
	for _, v := range inputMap {
		items = append(items, v)
	}

	sort.Slice(items, func(i, j int) bool {
		countI := strings.Count(items[i].Path, "/")
		countJ := strings.Count(items[j].Path, "/")

		if countI != countJ {
			return countI > countJ
		}

		return items[i].Path < items[j].Path
	})

	return items
}

// maxValue returns the maximum of two integers
func maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// isValidCustomTypeName checks if a given string is a valid custom type name
func isValidCustomTypeName(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Check if it's a standard type
	if stdTypes[s] {
		return false
	}

	// Check that the first character is a letter (considering Unicode) or underscore
	firstChar := rune(s[0])
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		return false
	}

	// Check the remaining characters
	for _, char := range s[1:] {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}

	return true
}

func deepCopy(src *Structure) *Structure {
	elem := &Structure{
		Name:        src.Name,
		Path:        src.Path,
		Root:        src.Root,
		RootField:   src.RootField,
		StructType:  src.StructType,
		StringType:  src.StringType,
		IsStructure: src.IsStructure,
		Size:        src.Size,
		Align:       src.Align,
		Offset:      src.Offset,
	}
	if src.MetaData != nil {
		elem.MetaData = &MetaData{
			BeforeSize: src.MetaData.BeforeSize,
			AfterSize:  src.MetaData.AfterSize,
			Data:       src.MetaData.Data,
			StartPos:   src.MetaData.StartPos,
			EndPos:     src.MetaData.EndPos,
		}
	}
	if src.NestedFields != nil {
		newNestedFields := make([]*Structure, 0, len(src.NestedFields))
		for _, item := range src.NestedFields {
			newNestedFields = append(newNestedFields, deepCopy(item))
		}
		elem.NestedFields = newNestedFields
	}

	return elem
}
