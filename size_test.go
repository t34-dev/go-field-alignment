package main

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"
	"unsafe"
)

func TestGetFieldSizeStruct(t *testing.T) {
	structExpr := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{Type: &ast.Ident{Name: "int"}},
				{Type: &ast.Ident{Name: "string"}},
				{Type: &ast.Ident{Name: "bool"}},
			},
		},
	}

	var expr ast.Expr = structExpr
	size := getFieldSize(expr)

	// Создаем реальную структуру для сравнения
	type testStruct struct {
		i int
		s string
		b bool
	}
	expected := reflect.TypeOf(testStruct{}).Size()

	// Выводим детальную информацию о размерах и выравнивании
	t.Logf("Размер int: %d", unsafe.Sizeof(int(0)))
	t.Logf("Размер string: %d", unsafe.Sizeof(""))
	t.Logf("Размер bool: %d", unsafe.Sizeof(false))
	t.Logf("Выравнивание структуры: %d", reflect.TypeOf(testStruct{}).Align())

	v := reflect.ValueOf(testStruct{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		t.Logf("Поле %s: размер %d, смещение %d, выравнивание %d",
			field.Name, field.Type.Size(), field.Offset, field.Type.Align())
	}

	t.Logf("Ожидаемый размер: %d", expected)
	t.Logf("Фактический размер: %d", size)

	if size != expected {
		t.Errorf("getFieldSize() for struct = %v, want %v", size, expected)
	}
}

func TestGetFieldSizeTypes(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{"bool", &ast.Ident{Name: "bool"}, unsafe.Sizeof(bool(false))},
		{"int", &ast.Ident{Name: "int"}, unsafe.Sizeof(int(0))},
		{"int8", &ast.Ident{Name: "int8"}, unsafe.Sizeof(int8(0))},
		{"int16", &ast.Ident{Name: "int16"}, unsafe.Sizeof(int16(0))},
		{"int32", &ast.Ident{Name: "int32"}, unsafe.Sizeof(int32(0))},
		{"int64", &ast.Ident{Name: "int64"}, unsafe.Sizeof(int64(0))},
		{"uint", &ast.Ident{Name: "uint"}, unsafe.Sizeof(uint(0))},
		{"uint8", &ast.Ident{Name: "uint8"}, unsafe.Sizeof(uint8(0))},
		{"uint16", &ast.Ident{Name: "uint16"}, unsafe.Sizeof(uint16(0))},
		{"uint32", &ast.Ident{Name: "uint32"}, unsafe.Sizeof(uint32(0))},
		{"uint64", &ast.Ident{Name: "uint64"}, unsafe.Sizeof(uint64(0))},
		{"float32", &ast.Ident{Name: "float32"}, unsafe.Sizeof(float32(0))},
		{"float64", &ast.Ident{Name: "float64"}, unsafe.Sizeof(float64(0))},
		{"complex64", &ast.Ident{Name: "complex64"}, unsafe.Sizeof(complex64(0))},
		{"complex128", &ast.Ident{Name: "complex128"}, unsafe.Sizeof(complex128(0))},
		{"string", &ast.Ident{Name: "string"}, unsafe.Sizeof("")},
		{"slice", &ast.ArrayType{Elt: &ast.Ident{Name: "int"}}, unsafe.Sizeof([]int{})},
		{"array", &ast.ArrayType{Elt: &ast.Ident{Name: "int"}, Len: &ast.BasicLit{Kind: token.INT, Value: "5"}}, unsafe.Sizeof([5]int{})},
		{"map", &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}}, unsafe.Sizeof(map[string]int{})},
		{"chan", &ast.ChanType{Value: &ast.Ident{Name: "int"}}, unsafe.Sizeof(make(chan int))},
		{"interface", &ast.InterfaceType{}, unsafe.Sizeof((*interface{})(nil))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFieldSize(tt.expr)
			if got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldSizePointers(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{"*int", &ast.StarExpr{X: &ast.Ident{Name: "int"}}, unsafe.Sizeof((*int)(nil))},
		{"*string", &ast.StarExpr{X: &ast.Ident{Name: "string"}}, unsafe.Sizeof((*string)(nil))},
		{"*bool", &ast.StarExpr{X: &ast.Ident{Name: "bool"}}, unsafe.Sizeof((*bool)(nil))},
		{"**int", &ast.StarExpr{X: &ast.StarExpr{X: &ast.Ident{Name: "int"}}}, unsafe.Sizeof((**int)(nil))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldSize(tt.expr); got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldSizeComplexStructs(t *testing.T) {
	tests := []struct {
		name string
		expr ast.Expr
		want uintptr
	}{
		{
			name: "struct with embedded struct",
			expr: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "int"}},
						{Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{Type: &ast.Ident{Name: "string"}},
									{Type: &ast.Ident{Name: "bool"}},
								},
							},
						}},
					},
				},
			},
			want: reflect.TypeOf(struct {
				i int
				s struct {
					str string
					b   bool
				}
			}{}).Size(),
		},
		{
			name: "struct with slice and map",
			expr: &ast.StructType{
				Fields: &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.ArrayType{Elt: &ast.Ident{Name: "int"}}},
						{Type: &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}}},
					},
				},
			},
			want: reflect.TypeOf(struct {
				s []int
				m map[string]int
			}{}).Size(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFieldSize(tt.expr); got != tt.want {
				t.Errorf("getFieldSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFieldSizeRecursiveStruct(t *testing.T) {
	// Определяем рекурсивную структуру
	type RecursiveStruct struct {
		i int
		r *RecursiveStruct
	}

	// Создаем AST представление структуры RecursiveStruct
	recursiveStructExpr := &ast.StructType{
		Fields: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{{Name: "i"}},
					Type:  &ast.Ident{Name: "int"},
				},
				{
					Names: []*ast.Ident{{Name: "r"}},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "RecursiveStruct"},
					},
				},
			},
		},
	}

	// Создаем полное AST представление определения типа
	typeSpec := &ast.TypeSpec{
		Name: &ast.Ident{Name: "RecursiveStruct"},
		Type: recursiveStructExpr,
	}

	// Оборачиваем определение типа в GenDecl, как это было бы в реальном Go файле
	genDecl := &ast.GenDecl{
		Tok:   token.TYPE,
		Specs: []ast.Spec{typeSpec},
	}

	// Теперь у нас есть полное AST представление определения типа RecursiveStruct

	var expr ast.Expr = recursiveStructExpr
	size := getFieldSize(expr)

	expected := reflect.TypeOf(RecursiveStruct{}).Size()

	t.Logf("AST представление структуры:")
	t.Logf("%#v", genDecl)
	t.Logf("Размер int: %d", unsafe.Sizeof(int(0)))
	t.Logf("Размер указателя: %d", unsafe.Sizeof(uintptr(0)))
	t.Logf("Ожидаемый размер: %d", expected)
	t.Logf("Фактический размер: %d", size)

	if size != expected {
		t.Errorf("getFieldSize() for recursive struct = %v, want %v", size, expected)
	}
}
