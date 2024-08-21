package main

import (
	"fmt"
	"strings"
)

// ============= Print
func viewPrintStructures(structures []*ItemInfo) {
	for _, elem := range structures {
		testPrintStructure(elem, 0)
		fmt.Printf("-------------------------------------------")
	}
}
func testPrintStructure(elem *ItemInfo, tab int) {
	// выравнивание для красивого отображения в логах
	maxFieldNameLength := 0
	maxTypeLength := 0

	if elem.IsStructure {
		for _, field := range elem.NestedFields {
			if len(field.Name) > maxFieldNameLength {
				maxFieldNameLength = len(field.Name)
			}
			if len(field.StringType) > maxTypeLength {
				maxTypeLength = len(field.StringType)
			}
		}
		infoFormat := fmt.Sprintf("%s    %%-%ds %%-%ds %%s", strings.Repeat(" ", tab), maxValue(maxFieldNameLength, 5), maxValue(maxTypeLength, 11))

		if tab == 0 {
			fmt.Printf("%stype %s struct {\n", strings.Repeat(" ", tab), elem.Name)
		} else {
			fmt.Printf("%s%s struct {\n", strings.Repeat(" ", tab), elem.Name)
		}
		var lastEnd uintptr
		for idx, field := range elem.NestedFields {
			isValidCustomNameType := isValidCustomTypeName(field.StringType)

			if field.IsStructure && !isValidCustomNameType {
				testPrintStructure(field, tab+4)
			} else {
				str := fmt.Sprintf("[Size: %d, Align: %d, Offset: %d]", field.Size, field.Align, field.Offset)
				padding := field.Offset - lastEnd
				if padding > 0 {
					str = fmt.Sprintf("+%db %s", padding, str)
				}
				lastEnd = field.Offset + field.Size
				if idx == len(elem.NestedFields)-1 {
					finalPadding := elem.Size - lastEnd
					if finalPadding > 0 {
						str = fmt.Sprintf("%s +%db", str, finalPadding)
					}
				}
				fmt.Printf(infoFormat+"\n", field.Name, field.StringType, str)
			}
		}
	}

	fmt.Printf("%s} [Size: %d, Align: %d]\n", strings.Repeat(" ", tab), elem.Size, elem.Align)
}
