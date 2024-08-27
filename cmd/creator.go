package main

import "go/ast"

// ============= Creator Item

// createItemInfoName extracts the name from the given AST node.
// It handles TypeSpec and Field nodes, returning an empty string for unsupported types.
func createItemInfoName(data interface{}) string {
	switch elem := data.(type) {
	case *ast.TypeSpec:
		return elem.Name.Name
	case *ast.Field:
		if len(elem.Names) == 0 {
			return ""
		}
		return elem.Names[0].Name
	}
	return ""
}

// createItemInfoPath generates a path string for an item.
// It combines the item's name with its parent's name, if available.
func createItemInfoPath(name, parentName string) string {
	if parentName != "" {
		return parentName + "/" + name
	}
	return name
}

// createItemInfo creates an Structure structure from the given AST node.
// It handles TypeSpec and Field nodes, processing their contents and creating nested structures as needed.
// The function also updates the provided mapper with the created Structure.
func createItemInfo(data interface{}, parentData *Structure, mapper map[string]*Structure) *Structure {
	switch Elem := data.(type) {
	case *ast.TypeSpec:
		newItem := &Structure{
			Name:        createItemInfoName(Elem),
			Root:        Elem,
			StructType:  Elem.Type,
			IsStructure: true,
			StringType:  getTypeString(Elem.Type),
		}
		if newItem.Name == "" {
			return nil
		}
		newItem.Path = createItemInfoPath(newItem.Name, "")
		mapper[newItem.Path] = newItem
		if elem, ok := Elem.Type.(*ast.StructType); ok {
			for _, field := range elem.Fields.List {
				newField := createItemInfo(field, newItem, mapper)
				if newField != nil {
					newItem.NestedFields = append(newItem.NestedFields, newField)
				}
			}
		}
		return newItem
	case *ast.Field:
		newItem := &Structure{
			Name:       createItemInfoName(Elem),
			RootField:  Elem,
			StructType: Elem.Type,
			StringType: getTypeString(Elem.Type),
		}
		if newItem.Name == "" {
			return nil
		}
		newItem.Path = createItemInfoPath(newItem.Name, parentData.Path)
		mapper[newItem.Path] = newItem
		if ident, ok := Elem.Type.(*ast.Ident); ok {
			newItem.IsStructure = ident.Obj != nil
		}
		if elem, ok := Elem.Type.(*ast.StructType); ok {
			newItem.IsStructure = true
			for _, field := range elem.Fields.List {
				newField := createItemInfo(field, newItem, mapper)
				if newField != nil {
					newItem.NestedFields = append(newItem.NestedFields, newField)
				}
			}
		}
		return newItem
	}
	return nil
}
