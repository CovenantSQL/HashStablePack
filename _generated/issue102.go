package _generated

//go:generate hsp

type Issue102 struct{}

type Issue102deep struct {
	A int
	X struct{}
	Y struct{}
	Z int
}

//hsp:tuple Issue102Tuple

type Issue102Tuple struct{}

//hsp:tuple Issue102TupleDeep

type Issue102TupleDeep struct {
	A int
	X struct{}
	Y struct{}
	Z int
}

type Issue102Uses struct {
	Nested    Issue102
	NestedPtr *Issue102
}

//hsp:tuple Issue102TupleUsesTuple

type Issue102TupleUsesTuple struct {
	Nested    Issue102Tuple
	NestedPtr *Issue102Tuple
}

//hsp:tuple Issue102TupleUsesMap

type Issue102TupleUsesMap struct {
	Nested    Issue102
	NestedPtr *Issue102
}

type Issue102MapUsesTuple struct {
	Nested    Issue102Tuple
	NestedPtr *Issue102Tuple
}
