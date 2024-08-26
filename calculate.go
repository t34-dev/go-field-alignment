package main

// calculateStructures calculates the size and alignment of structures in the given slice of Structure.
// It recursively processes nested structures and updates their size and alignment information.
func calculateStructure(elem *Structure, cache map[string]*Structure) {
	var currentOffset, maxAlign uintptr
	for _, field := range elem.NestedFields {
		var fieldSize, fieldAlign uintptr

		isValidCustomType := isValidCustomTypeName(field.StringType)

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
			if field.IsStructure {
				fieldSize, fieldAlign = calculateStructLayout(field)
			} else {
				fieldSize = getFieldSize(field.StructType)
				fieldAlign = getFieldAlign(field.StructType)
			}
		}

		currentOffset = align(currentOffset, fieldAlign)

		field.Size = fieldSize
		field.Align = fieldAlign
		field.Offset = currentOffset

		if isValidCustomType {
			cache[field.StringType] = field
		} else {
			cache[field.Path] = field
		}

		currentOffset += fieldSize

		if fieldAlign > maxAlign {
			maxAlign = fieldAlign
		}
	}

	elem.Size = align(currentOffset, maxAlign)
	elem.Align = maxAlign
}

func calculateStructLayout(field *Structure) (size, alignment uintptr) {
	var offset uintptr = 0
	maxAlign := uintptr(1)

	for _, field := range field.NestedFields {
		offset = align(offset, field.Align)
		if field.Align > maxAlign {
			maxAlign = field.Align
		}
		offset += field.Size
	}
	size = align(offset, maxAlign)
	alignment = maxAlign

	return size, alignment
}

// calculateStructure calculates the size and alignment of a single structure.
// It updates the Size and Align fields of the Structure and processes nested fields.
func calculateStructures(structures []*Structure, isBefore bool) {
	cache := make(map[string]*Structure, len(structures))
	for _, structure := range structures {
		calculateStructure(structure, cache)
		if isBefore {
			structure.MetaData.BeforeSize = structure.Size
		} else {
			structure.MetaData.AfterSize = structure.Size
		}
	}
}
