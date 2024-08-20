package main

import (
	"sort"
	"strings"
)

// ============= Optimization
func optimizeStructure(fields []*ItemInfo) []*ItemInfo {
	// Сортируем поля по убыванию выравнивания, затем по убыванию размера
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Align != fields[j].Align {
			return fields[i].Align > fields[j].Align
		}
		return fields[i].Size > fields[j].Size
	})

	// Отдельно обрабатываем массивы и слайсы
	var regularFields, arrayFields []*ItemInfo
	for _, field := range fields {
		if strings.HasPrefix(field.StringType, "[") || strings.HasPrefix(field.StringType, "[]") {
			arrayFields = append(arrayFields, field)
		} else {
			regularFields = append(regularFields, field)
		}
	}

	// Объединяем обратно, помещая массивы и слайсы в конец
	optimizedFields := append(regularFields, arrayFields...)

	// Пересчитываем смещения
	var currentOffset uintptr
	for i := range optimizedFields {
		currentOffset = align(currentOffset, optimizedFields[i].Align)
		optimizedFields[i].Offset = currentOffset
		currentOffset += optimizedFields[i].Size
	}

	return optimizedFields
}
func optimizeStructures(mapStructures map[string]*ItemInfo) {
	mapperItemsFlat := sortMapKeysBySlashCount(mapStructures)
	for _, structure := range mapperItemsFlat {
		if structure.IsStructure {
			structure.NestedFields = optimizeStructure(structure.NestedFields)
		}
	}
}
