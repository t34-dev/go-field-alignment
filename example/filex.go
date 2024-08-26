package example

// BadStructX is structure
type BadStructX struct {
	a  bool // 1 byte
	x1 struct {
		a  bool  // 1 byte
		b  bool  // 1 byte
		c  bool  // 1 byte
		c1 int32 // 1 byte
		c2 bool  // 1 byte
	}
	b  int32 // 4 bytes
	c  bool  // 1 byte
	d  int64 // 8 bytes
	xx struct {
		a bool  // 1 byte
		b int32 // 4 bytes
		c struct {
			a bool  // 1 byte
			b int32 `json:"logo" db:"logo" example:"http://url"` // 4 bytes
			c bool  // 1 byte
			d int64 // 8 bytes
		}
		d int64 // 8 bytes
	}
	x bool // 1 byte
} // BIG
