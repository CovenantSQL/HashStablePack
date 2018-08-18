package gen

import (
	"io"
	"text/template"
)

var (
	marshalTestTempl = template.New("MarshalTest")
	encodeTestTempl  = template.New("EncodeTest")
)

// TODO(philhofer):
// for simplicity's sake, right now
// we can only generate tests for types
// that can be initialized with the
// "Type{}" syntax.
// we should support all the types.

func mtest(w io.Writer) *mtestGen {
	return &mtestGen{w: w}
}

type mtestGen struct {
	passes
	w io.Writer
}

func (m *mtestGen) Execute(p Elem) error {
	p = m.applyall(p)
	if p != nil && IsPrintable(p) {
		switch p.(type) {
		case *Struct, *Array, *Slice, *Map:
			return marshalTestTempl.Execute(m.w, p)
		}
	}
	return nil
}

func (m *mtestGen) Method() Method { return marshaltest }

//type etestGen struct {
//	passes
//	w io.Writer
//}
//
//func etest(w io.Writer) *etestGen {
//	return &etestGen{w: w}
//}
//
//func (e *etestGen) Execute(p Elem) error {
//	p = e.applyall(p)
//	if p != nil && IsPrintable(p) {
//		switch p.(type) {
//		case *Struct, *Array, *Slice, *Map:
//			return encodeTestTempl.Execute(e.w, p)
//		}
//	}
//	return nil
//}
//
//func (e *etestGen) Method() Method { return encodetest }

func init() {
	template.Must(marshalTestTempl.Parse(`func TestMarshalHash{{.TypeName}}(t *testing.T) {
	v := {{.TypeName}}{}
	binary.Read(rand.Reader, binary.BigEndian, &v)
	bts1, err := v.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	bts2, err := v.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}
}

func BenchmarkMarshalHash{{.TypeName}}(b *testing.B) {
	v := {{.TypeName}}{}
	b.ReportAllocs()
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		v.MarshalHash()
	}
}

func BenchmarkAppendMsg{{.TypeName}}(b *testing.B) {
	v := {{.TypeName}}{}
	bts := make([]byte, 0, v.Msgsize())
	bts, _ = v.MarshalHash()
	b.SetBytes(int64(len(bts)))
	b.ReportAllocs()
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		bts, _ = v.MarshalHash()
	}
}

`))

}
