package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"sort"
	"strings"
	"unicode"
)

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

func formatGoCode(code string) (string, error) {
	formatted, err := format.Source([]byte(code))
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

// Helper function to get a string representation of the type
func getTypeString(expr ast.Expr) string {
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
	case *ast.StructType:
		return "struct{...}"
	default:
		return fmt.Sprintf("%T", expr)
	}
}
func sortMapKeysBySlashCount(inputMap map[string]*ItemInfo) []*ItemInfo {
	items := make([]*ItemInfo, 0, len(inputMap))
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

// Вспомогательная функция для определения максимума
func maxValue(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func isValidCustomTypeName(s string) bool {
	if len(s) == 0 {
		return false
	}

	// Проверяем, не является ли это стандартным типом
	if stdTypes[s] {
		return false
	}

	// Проверяем, что первый символ - буква (учитывая Unicode) или подчеркивание
	firstChar := rune(s[0])
	if !unicode.IsLetter(firstChar) && firstChar != '_' {
		return false
	}

	// Проверяем остальные символы
	for _, char := range s[1:] {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}

	return true
}
