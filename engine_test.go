package main

import (
	"fmt"
	textreplacer "github.com/t34-dev/go-text-replacer"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

const structsSourceIn = `
package main

//type BadStruct struct {
//	arr []string 
//	arr2 [5]int8 
//	a bool  // 1 byte
//	a1 string  // 1 byte
//	a2 *string  // 1 byte
//	b int32 // 4 bytes
//	c bool  // 1 byte
//	d int64 // 8 bytes
//	arr []string 
//	arr1 [3]string 
//	arr2 [3]int8 
//	arr3 [3]int8
//	arr4 [3]int8
//}

//type GoodStruct struct {
//	d int64 // 8 bytes
//	b int32 // 4 bytes
//	a bool  // 1 byte
//	c bool  // 1 byte
//}

// Test comment
type MyTest struct {
	a bool  // 1 byte
    nameX string
	b bool  // 1 byte
    App   struct {
// LogLevel
        LogLevel                  string        ` + "`yaml:\"log_level\" env-default:\"info\"` // 2 text" + `
        Name                      string        ` + "`yaml:\"name\" env-default:\"ms-sso\"`" + `
        IsProduction              bool          ` + "`yaml:\"is_production\" env:\"IS_PRODUCTION\" yaml-default:\"true\"`" + `
        TimeToConfirmRegistration time.Duration ` + "`yaml:\"tim_to_confirm_registration\" env-required:\"24h\"`" + `
    } ` + "`yaml:\"app\"`" + `
} // some text

//type MixStruct struct {
//	// Multi-line content
//	// Multi-line text
//	c bool  // 1 byte
//	a1 MyTest
//	a2 *MyTest
//	a3 struct {
//		n1 string
//		n2 struct {
//			// Some text
//			n0 bool
//			n1 string
//			n2 bool
//		}
//	}
//}
`
const structsSourceOut = `
package main

type BadStruct struct {
	a bool
	a1 string
	a2 *string
	b int32
	c bool
	d int64
	arr []string
	arr1 [3]string
	arr2 [3]int8
	arr3 [3]int8
	arr4 [3]int8
}

type GoodStruct struct {
	d int64
	b int32
	a bool
	c bool
}

type MixStruct struct {
	a1 BadStruct
	a3 struct {
		n1 string
		n2 struct {
			n1 string
			n2 int
		}
	}
	a2 GoodStruct
}
`

// TestStructAlignment tests the alignment and optimization of struct fields
func TestStructAlignment(t *testing.T) {
	fset := token.NewFileSet()

	// Parse the source code from the string
	node, err := parser.ParseFile(fset, "", structsSourceIn, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	var items []*ItemInfo
	mapperItems := map[string]*ItemInfo{}

	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			item := createItemInfo(typeSpec, nil, mapperItems)
			if item != nil {
				items = append(items, item)
			}
		}
		return true
	})
	calculateStructures(items)
	testPrintStructures(t, items)

	// Optimize structures
	optimizeStructures(mapperItems)
	calculateStructures(items)
	testPrintStructures(t, items)

	fmt.Println("===============")

	for _, structure := range items {
		data := renderStructure(structure)
		code, err := formatGoCode(data)
		if err != nil {
			fmt.Println("ERROR")
		}
		fmt.Println("=============== " + structure.Name)
		fmt.Println(code)
		fmt.Println("=============== END: " + structure.Name)
	}
	fmt.Println("items", items, mapperItems)

	// Replace content
	var blocks []textreplacer.Block
	for _, elem := range items {
		blocks = append(blocks, textreplacer.Block{
			Start: 1,
			End:   2,
			Txt:   []byte(elem.Name),
		})
	}

	replacer := textreplacer.New([]byte(structsSourceIn))
	result, err := replacer.Enter(blocks)

	modifiedCode := string(result)
	// Compare modified code with structsSourceOut
	expectedCode := strings.TrimSpace(structsSourceOut)
	if strings.TrimSpace(modifiedCode) != expectedCode {
		t.Errorf("Modified code does not match expected output.\nGot:\n%s\nWant:\n%s", modifiedCode, expectedCode)
	} else {
		t.Log("Modified code matches expected output.")
	}
}
