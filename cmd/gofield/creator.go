package main

import (
	"go/ast"
)

// ============= Creator Item

// createTypeName extracts the name from a type spec.
func createTypeName(spec *ast.TypeSpec) string {
	return spec.Name.Name
}

// createFieldNames extracts field names from a field spec.
//
// Returns multiple names when the specified field spec is actually a definition of multiple same-typed fields.
// Returns a slice with a single empty string in case of embedded types.
//
//		type A1 struct{}
//
//		type A2 struct {
//		  A1                // <-------- []string{""}
//		  F1 int            // <-------- []string{"F1"}
//		  F2 string         // <-------- []string{"F2"}
//	      F3, F4, F5 string // <-------- []string{"F3", "F4", "F5"}
//		}
func createFieldNames(field *ast.Field) []string {
	c := len(field.Names)
	if c == 0 {
		// This is a case of an embedded type - return a slice with an empty string
		return []string{""}
	}
	names := make([]string, c)
	for i := 0; i < c; i++ {
		names[i] = field.Names[i].Name
	}
	return names
}

// createItemInfoPath generates a path string for an item.
// It combines the item's name with its parent's name, if available.
func createItemInfoPath(name, parentName string) string {
	if parentName != "" {
		return parentName + "/" + name
	}
	return name
}

// createTypeItemInfo creates a Structure from the given AST type node.
// It processes its contents and creates nested structures as needed.
// The function also updates the provided mapper with the created Structure.
func createTypeItemInfo(typeSpec *ast.TypeSpec, parent *Structure, mapper map[string]*Structure) *Structure {
	typeInfo := &Structure{
		Name:        createTypeName(typeSpec),
		Root:        typeSpec,
		StructType:  typeSpec.Type,
		IsStructure: true,
		StringType:  getTypeString(typeSpec.Type),
	}
	if typeInfo.Name == "" {
		typeInfo.Name = "!" + typeInfo.StringType
	}

	typeInfo.Path = createItemInfoPath(typeInfo.Name, "")
	mapper[typeInfo.Path] = typeInfo

	if structType, ok := typeSpec.Type.(*ast.StructType); ok {
		for _, field := range structType.Fields.List {
			newFields := createFieldItemsInfo(field, typeInfo, mapper)
			if len(newFields) > 0 {
				typeInfo.NestedFields = append(typeInfo.NestedFields, newFields...)
			}
		}
	}

	return typeInfo
}

// createFieldItemsInfo creates a list of Structure-s from the given AST field node.
// The number of returned Structure-s is defined by how field node looks like.
// It processes its contents and creates nested structures as needed.
// The function also updates the provided mapper with the created Structures.
func createFieldItemsInfo(field *ast.Field, parent *Structure, mapper map[string]*Structure) []*Structure {
	fieldNames := createFieldNames(field)
	fieldInfos := make([]*Structure, len(fieldNames))
	for i, name := range fieldNames {
		fieldInfos[i] = createSingleFieldItemInfo(name, field, parent, mapper)
	}
	return fieldInfos
}

// createSingleFieldItemInfo creates a single Structure from the given AST field node.
// It processes its contents and creates nested structures as needed.
// The function also updates the provided mapper with the created Structure.
func createSingleFieldItemInfo(name string, field *ast.Field, parent *Structure, mapper map[string]*Structure) *Structure {
	fieldInfo := &Structure{
		Name:       name,
		RootField:  field,
		StructType: field.Type,
		StringType: getTypeString(field.Type),
	}
	if fieldInfo.Name == "" {
		fieldInfo.Name = "!" + fieldInfo.StringType
	}

	fieldInfo.Path = createItemInfoPath(fieldInfo.Name, parent.Path)
	mapper[fieldInfo.Path] = fieldInfo

	switch typed := field.Type.(type) {
	case *ast.Ident:
		fieldInfo.IsStructure = typed.Obj != nil
	case *ast.StructType:
		fieldInfo.IsStructure = true
		for _, nestedField := range typed.Fields.List {
			nestedFieldItems := createFieldItemsInfo(nestedField, fieldInfo, mapper)
			if len(nestedFieldItems) > 0 {
				fieldInfo.NestedFields = append(fieldInfo.NestedFields, nestedFieldItems...)
			}
		}
	}

	return fieldInfo
}
