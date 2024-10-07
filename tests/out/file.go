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

type (
	S1 struct {
		F2 string
		F1 bool
	}

	S2 struct {
		F3 string
		F4 StructWithGenerics[int]
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
