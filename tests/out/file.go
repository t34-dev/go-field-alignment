package enter

import (
	"fmt"
	"time"
)

type Problem1 struct{}
type Problem2 struct{}
type Problem3 struct {
	hello  string
	hello2 string
	time.Time
	time.Duration
	time.Location
	typer bool
}

type StructWithGenerics[T any] struct {
	F T
}

type StructWithMoreGenerics[T1, T2 any, T3 comparable] struct {
	F1 T1
	F2 T2
	F3 T3
}

type (
	S1 struct {
		F2 string
		F1 bool
	}

	S2 struct {
		F3 string
		F4 StructWithGenerics[int]
		F5 StructWithMoreGenerics[int, float64, string]
		F1 bool
		F2 bool
	}
)

type STR string
type STRs []string

// Test comment
type MyTest struct {
	App struct {
		// LogLevel
		LogLevel                  string        `yaml:"log_level" env-default:"info"` // 2 text
		Name                      string        `yaml:"name" env-default:"ms-sso"`
		TimeToConfirmRegistration time.Duration `yaml:"tim_to_confirm_registration" env-required:"24h"`
		IsProduction              bool          `yaml:"is_production" env:"IS_PRODUCTION" yaml-default:"true"`
	} `yaml:"app"`
	nameX    string
	Problem1 struct {
		I interface{}
		S struct{}
	}
	a bool // 1 byte
	b bool // 1 byte
} /* some text
dsdsd
dsds
*/

type MyTest2 struct{}

type MyTest3 struct{}

type MyTest4 struct{} // test

type MyTest5 struct{} // test

var Name = "dsds"

func Get() {
	f := MyTest{}
	f2 := MyTest2{}
	f3 := MyTest3{}
	f4 := MyTest4{}
	f5 := MyTest5{}
	fmt.Println(f, f2, f3, f4, f5, Name)

	//name
}

// Unaligned Structures
type SimpleGenericUnaligned[T any] struct {
	Name  string
	ID    int
	Value T
}

type MultiParamUnaligned[T string, U Number] struct {
	IsValid bool
	First   T
	Second  U
}

type OuterUnaligned[T any, U comparable] struct {
	Nested   Inner[U]
	Priority int
	Data     T
}

type TreeNodeUnaligned[T any] struct {
	Left   *TreeNodeUnaligned[T]
	Right  *TreeNodeUnaligned[T]
	Depth  int
	IsLeaf bool
	Value  T
}

type SliceContainerUnaligned[T any] struct {
	TotalCount int
	MaxSize    int64
	IsReadOnly bool
	Items      []T
}

// Aligned Structures
type SimpleGenericAligned[T any] struct {
	Name  string
	ID    int
	Value T
}

type MultiParamAligned[T any, U Number] struct {
	IsValid bool
	First   T
	Second  U
}

type OuterAligned[T any, U comparable] struct {
	Nested   Inner[U]
	Priority int
	Data     T
}

type TreeNodeAligned[T any] struct {
	Left   *TreeNodeAligned[T]
	Right  *TreeNodeAligned[T]
	Depth  int
	IsLeaf bool
	Value  T
}

type SliceContainerAligned[T any] struct {
	TotalCount int
	MaxSize    int64
	IsReadOnly bool
	Items      []T
}

// Common types used in the structures above
type Number interface {
	int | int32 | int64 | float32 | float64
}

type Inner[T comparable] struct {
	Value string
	Key   T
}
