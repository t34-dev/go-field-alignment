package main

import (
	"fmt"
	"go/ast"
	"strings"
)

// ============= Render

// renderStructure generates a string representation of an Structure structure.
// It handles both top-level structures and nested fields, including their
// documentation, tags, and comments. The function recursively processes
// nested structures to create a complete representation.
//
// The function performs the following tasks:
// - Checks if the element is a valid custom type or a structure
// - Generates the struct definition with its name (for top-level structures)
// - Iterates through all nested fields, rendering each one
// - Includes field documentation, tags, and comments
// - Handles both root-level and nested comments
//
// Parameters:
// - elem: Pointer to an Structure structure to be rendered
//
// Returns:
// - A string containing the rendered structure
func renderStructure(elem *Structure) string {
	isValidCustomNameType := isValidCustomTypeName(elem.StringType)
	if !elem.IsStructure || isValidCustomNameType {
		return elem.StringType
	}

	var data strings.Builder
	if elem.Root != nil {
		data.WriteString(elem.Name)

		// Render generic type params
		if typeParams := elem.Root.TypeParams; typeParams != nil && len(typeParams.List) > 0 {
			params := typeParams.List
			data.WriteRune('[')
			data.WriteString(renderTypeParameter(params[0]))
			for _, p := range params[1:] {
				data.WriteString(", ")
				data.WriteString(renderTypeParameter(p))
			}
			data.WriteRune(']')
		}

		// Don't add "type " here, as current structure may be inside a "type" block
		data.WriteString(" struct {")
	} else {
		// Anonymous structs don't support type params
		data.WriteString("struct {")
	}
	if len(elem.NestedFields) > 0 {
		data.WriteRune('\n')
	}

	for idx, field := range elem.NestedFields {
		// Doc
		if field.RootField != nil && field.RootField.Doc != nil && len(field.RootField.Doc.List) > 0 {
			for _, comment := range field.RootField.Doc.List {
				data.WriteString(comment.Text)
				data.WriteRune('\n')
			}
		}
		if strings.HasPrefix(field.Name, "!") {
			field.Name = ""
		}
		data.WriteString(fmt.Sprintf("%s %s ", field.Name, renderStructure(field)))
		// Tag
		if field.RootField != nil {
			// Tags
			if field.RootField.Tag != nil && len(field.RootField.Tag.Value) > 0 {
				data.WriteString(fmt.Sprintf(" %s", field.RootField.Tag.Value))
			}
			// Comment
			if field.RootField.Comment != nil && len(field.RootField.Comment.List) > 0 {
				for _, comment := range field.RootField.Comment.List {
					data.WriteString(fmt.Sprintf(" %s", comment.Text))
				}
			}
		}
		if idx != len(elem.NestedFields) {
			data.WriteRune('\n')
		}
	}

	data.WriteRune('}')

	// Comments
	if elem.RootField != nil {
		if elem.RootField.Comment != nil && len(elem.RootField.Comment.List) > 0 {
			for _, comment := range elem.RootField.Comment.List {
				data.WriteString(comment.Text)
			}
		}
	} else if elem.Root != nil {
		if elem.Root.Comment != nil && len(elem.Root.Comment.List) > 0 {
			for _, comment := range elem.Root.Comment.List {
				data.WriteString(comment.Text)
			}
		}
	}
	return data.String()
}

// renderTypeParameter renders the given type parameter as Go code.
func renderTypeParameter(f *ast.Field) string {
	names := make([]string, len(f.Names))
	for i := 0; i < len(names); i++ {
		names[i] = f.Names[i].Name
	}
	return fmt.Sprintf("%s %s", strings.Join(names, ", "), getTypeString(f.Type))
}

func renderTextStructures(structures []*Structure) {
	for _, structure := range structures {
		// Don't format code here - "renderStructure" generates a replacement for a part of target Go file,
		// not a valid piece of Go code per-se.
		//
		// Code will be formatted afterwards.
		structure.MetaData.Data = []byte(renderStructure(structure))
	}
}
