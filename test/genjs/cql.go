package covenant

//go:generate hsp
type Person struct {
	Name   string
	Age    uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64
	Age2   int
	Int8   int8
	Int16  int16
	Int32  int32
	Int64  int64
	F1     float32
	F2     float64

	// Map         map[string]int
	unexported1 bool // this field is ignored
}
