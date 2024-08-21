package fileout

import (
	"fmt"
	"time"
)

// Test comment
type MyTest struct {
	a     bool // 1 byte
	nameX string
	b     bool // 1 byte
	App   struct {
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

type MyTest2 struct {
}

type MyTest3 struct {
}

type MyTest4 struct {
} // test

type MyTest5 struct {
} // test

var Name = "dsds"

func Get() {
	f := MyTest{}
	f2 := MyTest2{}
	f3 := MyTest3{}
	f4 := MyTest4{}
	f5 := MyTest5{}
	fmt.Println(f, f2, f3, f4, f5, Name)
}
