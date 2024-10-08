package enter

import (
	"fmt"
	"time"
)

type Problem1 struct{}
type Problem2 struct {
}
type Problem3 struct {
	hello, hello2 string
	typer         bool
	time.Time
	time.Duration
	time.Location
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
		F1 bool
		F2 string
	}

	S2 struct {
		F1, F2 bool
		F3     string
		F4     StructWithGenerics[int]
		F5     StructWithMoreGenerics[int, float64, string]
	}
)

type STR string
type STRs []string

// Test comment
type MyTest struct {
	a        bool // 1 byte
	nameX    string
	Problem1 struct {
		I interface{}
		S struct{}
	}
	b   bool // 1 byte
	App struct {
		// LogLevel
		LogLevel                  string        `yaml:"log_level" env-default:"info"` // 2 text
		Name                      string        `yaml:"name" env-default:"ms-sso"`
		IsProduction              bool          `yaml:"is_production" env:"IS_PRODUCTION" yaml-default:"true"`
		TimeToConfirmRegistration time.Duration `yaml:"tim_to_confirm_registration" env-required:"24h"`
	} `yaml:"app"`
} /* some text
dsdsd
dsds
*/

type MyTest2 struct{}

type MyTest3 struct {
}

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
	Value T
	ID    int
	Name  string
}

type MultiParamUnaligned[T string, U Number] struct {
	First   T
	Second  U
	IsValid bool
}

type OuterUnaligned[T any, U comparable] struct {
	Data     T
	Nested   Inner[U]
	Priority int
}

type TreeNodeUnaligned[T any] struct {
	Value  T
	Left   *TreeNodeUnaligned[T]
	Right  *TreeNodeUnaligned[T]
	Depth  int
	IsLeaf bool
}

type SliceContainerUnaligned[T any] struct {
	Items      []T
	TotalCount int
	MaxSize    int64
	IsReadOnly bool
}

// Aligned Structures
type SimpleGenericAligned[T any] struct {
	Value T
	ID    int
	Name  string
}

type MultiParamAligned[T any, U Number] struct {
	First   T
	Second  U
	IsValid bool
}

type OuterAligned[T any, U comparable] struct {
	Data     T
	Nested   Inner[U]
	Priority int
}

type TreeNodeAligned[T any] struct {
	Value  T
	Left   *TreeNodeAligned[T]
	Right  *TreeNodeAligned[T]
	Depth  int
	IsLeaf bool
}

type SliceContainerAligned[T any] struct {
	Items      []T
	TotalCount int
	MaxSize    int64
	IsReadOnly bool
}

// Common types used in the structures above
type Number interface {
	int | int32 | int64 | float32 | float64
}

type Inner[T comparable] struct {
	Key   T
	Value string
}
