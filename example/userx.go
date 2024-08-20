package example

type ExampleExp struct {
	E struct {
		D struct {
			B int64 `json:"logo" db:"logo" example:"http://url"`
			F int64
			D int32
			A bool `json:"id" db:"id"`
			C bool // тут текст

			E bool
		}
		B int64 `json:"logo" db:"logo" example:"http://url"`
		F int64 // тут коммент

		A bool `json:"id" db:"id"`
		C bool // тут текст

		E bool
	}
	B int64 `json:"logo" db:"logo" example:"http://url"`
	F int64
	D int32
	A bool `json:"id" db:"id"`
	C bool // тут текст

}
