package gen

import (
	"fmt"
	"io"
	"sort"

	"github.com/CovenantSQL/HashStablePack/marshalhash"
)

func marshal(w io.Writer) *marshalGen {
	return &marshalGen{
		p: printer{w: w},
	}
}

type marshalGen struct {
	passes
	p    printer
	fuse []byte
	v    string
}

func (m *marshalGen) Method() Method { return Marshal }

func (m *marshalGen) Apply(dirs []string) error {
	return nil
}

func (m *marshalGen) sort(e Elem) {
	if es, ok := e.(*Struct); ok {
		sort.Sort(es)
	}
}

func (m *marshalGen) setVersion(v string) {
	m.v = v
}

func (m *marshalGen) Execute(p Elem) error {
	if !m.p.ok() {
		return m.p.err
	}
	p = m.applyall(p)
	if p == nil {
		return nil
	}
	m.sort(p)

	if !IsPrintable(p) {
		return nil
	}

	m.p.comment("MarshalHash marshals for hash")

	// save the vname before
	// calling methodReceiver so
	// that z.Msgsize() is printed correctly
	c := p.Varname()

	if ps, ok := p.(*Struct); ok && ps.Versioning && m.v == "" {
		// print version info list
		m.p.printf("\nvar hspVersions%s = []string{", p.TypeName())
		for i := range ps.VersionList {
			m.p.printf("\n\"%s\",", ps.VersionList[i])
		}
		m.p.print("\n}")
		// print current version function
		m.p.printf("\nfunc (%s %s) HSPCurrentVersion() int {", c, imutMethodReceiver(p))
		m.p.printf("\nreturn int(%s.%s)", c, ps.VersionField)
		m.p.print("\n}")
		// print max version function
		m.p.printf("\nfunc (%s %s) HSPMaxVersion() int {", c, imutMethodReceiver(p))
		m.p.printf("\nreturn %d", len(ps.VersionList)-1)
		m.p.print("\n}")
		// print default version function
		m.p.printf("\nfunc (%s %s) HSPDefaultVersion() int {", c, imutMethodReceiver(p))
		m.p.printf("\nreturn %d", ps.CurrentNumericVersion)
		m.p.print("\n}")
	}

	m.p.printf("\nfunc (%s %s) MarshalHash%s() (o []byte, err error) ", c, imutMethodReceiver(p), m.v)
	if m.v != "oldver" {
		m.p.printf("{")

		if ps, ok := p.(*Struct); ok && ps.Versioning && m.v == "" {
			// version enabled and print switch statements
			m.p.printf("\nswitch %s.HSPCurrentVersion() {", c)
			for i := range ps.VersionList {
				m.p.printf("\ncase %d:", i)
				m.p.printf("\nreturn %s.MarshalHash%s()", c, ps.VersionList[i])
			}
			m.p.print("\ndefault:")
			m.p.print("\nerr = errors.New(\"invalid struct version\")")
			m.p.print("\nreturn")
			m.p.print("\n}")
			m.p.nakedReturn()
		} else {
			m.p.printf("\nvar b []byte")
			m.p.printf("\no = hsp.Require(b, %s.Msgsize%s())", c, m.v)
			next(m, p)
			m.p.nakedReturn()
		}
	} else {
		m.p.print(p.(*Struct).OldMarshalBody)
	}

	return m.p.err
}

func (m *marshalGen) rawAppend(typ string, argfmt string, arg interface{}) {
	m.p.printf("\no = hsp.Append%s(o, %s)", typ, fmt.Sprintf(argfmt, arg))
}

func (m *marshalGen) fuseHook() {
	if len(m.fuse) > 0 {
		m.rawbytes(m.fuse)
		m.fuse = m.fuse[:0]
	}
}

func (m *marshalGen) Fuse(b []byte) {
	if len(m.fuse) == 0 {
		m.fuse = b
	} else {
		m.fuse = append(m.fuse, b...)
	}
}

func (m *marshalGen) gStruct(s *Struct) {
	if !m.p.ok() {
		return
	}

	if s.AsTuple {
		m.tuple(s)
	} else {
		m.mapstruct(s)
	}
	return
}

func (m *marshalGen) tuple(s *Struct) {
	data := make([]byte, 0, 5)
	data = marshalhash.AppendArrayHeader(data, uint32(len(s.Fields)))
	m.p.printf("\n// array header, size %d", len(s.Fields))
	m.Fuse(data)
	if len(s.Fields) == 0 {
		m.fuseHook()
	}
	for i := range s.Fields {
		if !m.p.ok() {
			return
		}
		next(m, s.Fields[i].FieldElem)
	}
}

func (m *marshalGen) mapstruct(s *Struct) {
	data := make([]byte, 0, 64)
	data = marshalhash.AppendMapHeader(data, uint32(len(s.Fields)))
	m.p.printf("\n// map header, size %d", len(s.Fields))
	m.Fuse(data)
	if len(s.Fields) == 0 {
		m.fuseHook()
	}
	for i := range s.Fields {
		if !m.p.ok() {
			return
		}
		//data = hsp.AppendString(nil, s.Fields[i].FieldTag)
		//
		//m.p.printf("\n// string %q", s.Fields[i].FieldTag)
		//m.Fuse(data)

		next(m, s.Fields[i].FieldElem)
	}
}

// append raw data
func (m *marshalGen) rawbytes(bts []byte) {
	m.p.print("\no = append(o, ")
	for _, b := range bts {
		m.p.printf("0x%x,", b)
	}
	m.p.print(")")
}

func (m *marshalGen) gMap(s *Map) {
	if !m.p.ok() {
		return
	}
	m.fuseHook()
	vname := s.Varname()
	m.rawAppend(mapHeader, lenAsUint32, vname)
	m.p.printf("\n%sSlice := make([]string, 0, len(%s))", s.Keyidx, vname)
	m.p.printf("\nfor i := range %s {\n%sSlice = append(%sSlice, i)\n}",
		vname, s.Keyidx, s.Keyidx)
	m.p.printf("\nsort.Strings(%sSlice)", s.Keyidx)
	m.p.printf("\nfor _, %s := range %sSlice {\n %s := %s[%s]", s.Keyidx, s.Keyidx, s.Validx, vname, s.Keyidx)
	//m.p.printf("\nfor %s, %s := range %s {", s.Keyidx, s.Validx, vname)
	m.rawAppend(stringTyp, literalFmt, s.Keyidx)
	next(m, s.Value)
	m.p.closeblock()
}

func (m *marshalGen) gSlice(s *Slice) {
	if !m.p.ok() {
		return
	}
	m.fuseHook()
	vname := s.Varname()
	m.rawAppend(arrayHeader, lenAsUint32, vname)
	m.p.rangeBlock(s.Index, vname, m, s.Els)
}

func (m *marshalGen) gArray(a *Array) {
	if !m.p.ok() {
		return
	}
	m.fuseHook()
	if be, ok := a.Els.(*BaseElem); ok && be.Value == Byte {
		m.rawAppend("Bytes", "(%s)[:]", a.Varname())
		return
	}

	m.rawAppend(arrayHeader, literalFmt, coerceArraySize(a.Size))
	m.p.rangeBlock(a.Index, a.Varname(), m, a.Els)
}

func (m *marshalGen) gPtr(p *Ptr) {
	if !m.p.ok() {
		return
	}
	m.fuseHook()
	m.p.printf("\nif %s == nil {\no = hsp.AppendNil(o)\n} else {", p.Varname())
	next(m, p.Value)
	m.p.closeblock()
}

func (m *marshalGen) gBase(b *BaseElem) {
	if !m.p.ok() {
		return
	}
	m.fuseHook()
	vname := b.Varname()

	if b.Convert {
		if b.ShimMode == Cast {
			vname = tobaseConvert(b)
		} else {
			vname = randIdent()
			m.p.printf("\nvar %s %s", vname, b.BaseType())
			m.p.printf("\n%s, err = %s", vname, tobaseConvert(b))
			m.p.printf(errcheck)
		}
	}

	var echeck bool
	switch b.Value {
	case IDENT:
		m.p.printf(`
			if oTemp, err := %s.MarshalHash(); err != nil {
				return nil, err
			} else {
				o = hsp.AppendBytes(o, oTemp)
			}`, vname)
	case Intf, Ext:
		echeck = true
		m.p.printf("\no, err = hsp.Append%s(o, %s)", b.BaseName(), vname)
	default:
		m.rawAppend(b.BaseName(), literalFmt, vname)
	}

	if echeck {
		m.p.print(errcheck)
	}
}
