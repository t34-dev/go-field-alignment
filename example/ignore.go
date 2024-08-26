package example

type BadStruct struct {
	a bool  // 1 byte
	b int32 // 4 bytes
	c bool  // 1 byte
	d int64 // 8 bytes
} // test

// Title
type BadStruct2 struct {
	a bool  // 1 byte
	b int32 // 4 bytes
	c bool  // 1 byte
	d int64 // 8 bytes
}

/*
Long
Text
*/
type BadStruct3 struct {
	d int64 // 8 bytes
	b int32 // 4 bytes
	a bool  // 1 byte
	c bool  // 1 byte
}
type BadStruct4 struct {
	b int32 // 4 bytes
	d int64 // 8 bytes
	a bool  // 1 byte
	c bool  // 1 byte
}
