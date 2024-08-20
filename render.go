package main

import "fmt"

// ============= Render
func renderStructure(elem *ItemInfo) string {
	data := ""

	isValidCustomNameType := isValidCustomTypeName(elem.StringType)
	if !elem.IsStructure || isValidCustomNameType {
		return elem.StringType
	}

	if elem.Root != nil {
		data += fmt.Sprintf("type ")
		data += fmt.Sprintf("%s struct{\n", elem.Name)
	} else {
		data += fmt.Sprintf("struct {\n")
	}

	for idx, field := range elem.NestedFields {
		// Doc
		if field.RootField != nil && field.RootField.Doc != nil && len(field.RootField.Doc.List) > 0 {
			for _, comment := range field.RootField.Doc.List {
				data += fmt.Sprintln(comment.Text)
			}
		}
		data += fmt.Sprintf("%s %s ", field.Name, renderStructure(field))
		// Tag
		if field.RootField != nil {
			// Tags
			if field.RootField.Tag != nil && len(field.RootField.Tag.Value) > 0 {
				data += fmt.Sprintf("%s ", field.RootField.Tag.Value)
			}
			// Comment
			if field.RootField.Comment != nil && len(field.RootField.Comment.List) > 0 {
				for _, comment := range field.RootField.Comment.List {
					data += fmt.Sprintf("%s ", comment.Text)
				}
			}
		}
		if idx != len(elem.NestedFields) {
			data += fmt.Sprintf("\n")
		}
	}

	data += fmt.Sprintf("}")
	// Comments
	if elem.RootField != nil {
		if elem.RootField.Comment != nil && len(elem.RootField.Comment.List) > 0 {
			for _, comment := range elem.RootField.Comment.List {
				data += fmt.Sprintf("%s ", comment.Text)
			}
		}
	} else if elem.Root != nil {
		if elem.Root.Comment != nil && len(elem.Root.Comment.List) > 0 {
			for _, comment := range elem.Root.Comment.List {
				data += fmt.Sprintf("%s ", comment.Text)
			}
		}
	}
	return data
}
