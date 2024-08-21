package main

// calculateStructures calculates the size and alignment of structures in the given slice of ItemInfo.
// It recursively processes nested structures and updates their size and alignment information.
func calculateStructure(elem *ItemInfo, cache map[string]*ItemInfo) {
	var currentOffset, maxAlign uintptr
	for _, field := range elem.NestedFields {
		var fieldSize, fieldAlign uintptr

		isValidTypeName := isValidCustomTypeName(field.StringType)

		if field.IsStructure {
			calculateStructure(field, cache)
		}

		if item, ok := cache[field.StringType]; ok {
			fieldSize = item.Size
			fieldAlign = item.Align
		} else if item, ok = cache[elem.Path]; ok {
			fieldSize = item.Size
			fieldAlign = item.Align
		} else {
			fieldSize = getFieldSize(field.StructType)
			fieldAlign = getFieldAlign(field.StructType)
		}

		currentOffset = align(currentOffset, fieldAlign)

		field.Size = fieldSize
		field.Align = fieldAlign
		field.Offset = currentOffset

		if isValidTypeName {
			cache[field.StringType] = field
		} else {
			cache[elem.Path] = field
		}

		currentOffset += fieldSize

		if fieldAlign > maxAlign {
			maxAlign = fieldAlign
		}
	}

	elem.Size = align(currentOffset, maxAlign)
	elem.Align = maxAlign
}

// calculateStructure calculates the size and alignment of a single structure.
// It updates the Size and Align fields of the ItemInfo and processes nested fields.
func calculateStructures(structures []*ItemInfo) {
	cache := make(map[string]*ItemInfo, len(structures))
	for _, structure := range structures {
		calculateStructure(structure, cache)
	}
}
