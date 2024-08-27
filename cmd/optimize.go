package main

import (
	"sort"
	"strings"
)

// ============= Optimization

// optimizeStructure reorganizes the fields of a structure to minimize padding and optimize memory usage.
// It sorts fields by alignment and size, separates regular fields from arrays and slices,
// and recalculates field offsets for the optimized structure.
func optimizeStructure(fields []*Structure) []*Structure {
	// Sort fields in descending order of alignment, then in descending order of size
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Align != fields[j].Align {
			return fields[i].Align > fields[j].Align
		}
		return fields[i].Size > fields[j].Size
	})

	// Separately process arrays and slices
	var regularFields, arrayFields []*Structure
	for _, field := range fields {
		if strings.HasPrefix(field.StringType, "[") || strings.HasPrefix(field.StringType, "[]") {
			arrayFields = append(arrayFields, field)
		} else {
			regularFields = append(regularFields, field)
		}
	}

	// Merge back, placing arrays and slices at the end
	optimizedFields := append(regularFields, arrayFields...)

	// Recalculate offsets
	var currentOffset uintptr
	for i := range optimizedFields {
		currentOffset = align(currentOffset, optimizedFields[i].Align)
		optimizedFields[i].Offset = currentOffset
		currentOffset += optimizedFields[i].Size
	}

	return optimizedFields
}

// optimizeMapperStructures applies the optimizeStructure function to all structures in the given map.
// It processes structures in order of their nesting depth (determined by the number of slashes in their path).
func optimizeMapperStructures(mapStructures map[string]*Structure) {
	mapperItemsFlat := sortMapKeysBySlashCount(mapStructures)
	for _, structure := range mapperItemsFlat {
		if structure.IsStructure {
			structure.NestedFields = optimizeStructure(structure.NestedFields)
		}
	}
}
