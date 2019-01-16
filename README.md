Hash Stable Pack
=======
This is a code generation tool and serialization library for Calculation of Stable Hash for content. Basically it will generate an `MarshalHash` method which follow the MessagePack but **without the key**. 

### For
- Quick compare nested struct without reflection
- Quick calculation of struct hash or signature without reflection

### How?

That is the following 2 structs with different member name
For more: see [Spec in Chinese](spec.md)
```go
package person

//go:generate hsp
type Person1 struct {
	Name1       string 
	Age1        int    
	Address1    string 
	unexported1 bool             // this field is ignored
}

// Same struct with "string, string, int, bool"
type Person2 struct {
	Name2       string 
	Address2    string 
	Age2        int    
	unexported2 bool             // this field is ignored
}
```

But with the same type and content of exported member, `MarshalHash` will produce the same bytes array:
```go
package person

import (
	"bytes"
	"testing"
)

func TestMarshalHashAccountStable3(t *testing.T) {
	p1 := Person1{
		Name1:       "Auxten",
		Age1:        28,
		Address1:    "@CovenantSQL.io",
		unexported1: false,
	}
	p2 := Person2{
		Name2:       "Auxten",
		Address2:    "@CovenantSQL.io",
		Age2:        28,
		unexported2: true,
	}
	bts1, err := p1.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	bts2, err := p2.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}
}
```
the order of struct member is sorted by type name, so "string, int", "int, string" is equivalent.



You can read more about MessagePack [in the wiki](http://github.com/tinylib/msgp/wiki), or at [msgpack.org](http://msgpack.org).

### Why?

- Use Go as your schema language
- Performance
- [JSON interop](http://godoc.org/github.com/tinylib/msgp/msgp#CopyToJSON)
- [User-defined extensions](http://github.com/tinylib/msgp/wiki/Using-Extensions)
- Type safety
- Encoding flexibility

### Why not?

- MessagePack: member name is unnecessary, different implementation may add some fields which made result undetermined.
- Prorobuf: struct must defined in proto language, and other limitations discussed [here](https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510)

### Quickstart

1. Quick Install
```bash
go get -u github.com/CovenantSQL/HashStablePack/hsp
```

2. Add tag for source
In a source file, include the following directive:

```go
//go:generate hsp
```

3. Run go generate
```bash
go generate ./...
```

The `hsp` command will generate serialization methods for all exported type declarations in the file.

By default, the code generator will only generate `MarshalHash` and `Msgsize` method
```go
func (z *Test) MarshalHash() (o []byte, err error)
func (z *Test) Msgsize() (s int)
```


### Features

 - Extremely fast generated code
 - Test and benchmark generation
 - Support for complex type declarations
 - Native support for Go's `time.Time`, `complex64`, and `complex128` types 
 - Support for arbitrary type system extensions
 - File-based dependency model means fast codegen regardless of source tree size.

Consider the following:
```go
const Eight = 8
type MyInt int
type Data []byte

type Struct struct {
	Which  map[string]*MyInt 
	Other  Data              
	Nums   [Eight]float64    
}
```
As long as the declarations of `MyInt` and `Data` are in the same file as `Struct`, the parser will determine that the type information for `MyInt` and `Data` can be passed into the definition of `Struct` before its methods are generated.

### Known issues
- map type is not supported. will cause undetermined marshal content.

### License

This lib is inspired by https://github.com/tinylib/msgp
Most Code is diverted from https://github.com/tinylib/msgp, but It's an total different lib for usage. So I created a new project instead of forking it.

