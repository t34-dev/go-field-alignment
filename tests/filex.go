package fileout

type BadStruct struct {
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
		c bool  // 1 byte
		d int64 // 8 bytes
	}
	x bool // 1 byte
}
