# HashStablePack Spec

HashStablePack（后面简称 HSP）是基于 `https://github.com/tinylib/msgp` 修改而来的。主要序列化策略遵循 MsgPack 的标准。

> 参见：https://github.com/msgpack/msgpack/blob/master/spec.md



主要改动如下：

1. 去掉了 MsgPack 的 Key
2. 对结构体的成员按照类型做了排序



## 排序策略

排序使用的上 golang 标准的 sort.Sort，比较函数如下：

```go
// Less returns true if node i is less than node j.
func (x *Struct) Less(i, j int) bool {
   fi := x.Fields[i]
   fj := x.Fields[j]
   return fi.FieldElem.TypeName() < fj.FieldElem.TypeName()
}
```

整体上是按照字段类型的名字进行排序，各种字段类型的名字`TypeName()`对应如下：

#### Array

```go
fmt.Sprintf("[%s]%s", a.Size, a.Els.TypeName())
```

#### Map

```go
"map[string]" + m.Value.TypeName()
```

#### 指针

```go
"*" + s.Value.TypeName()
```

#### Slice

```go
"[]" + s.Els.TypeName()
```

#### Struct

```go
"struct{}"
```

#### 基本类型

请参考这两个函数

```go
func (s *BaseElem) BaseType() string {
	switch s.Value {
	case IDENT:
		return s.TypeName()

	// exceptions to the naming/capitalization
	// rule:
	case Intf:
		return "interface{}"
	case Bytes:
		return "[]byte"
	case Time:
		return "time.Time"
	case Ext:
		return "hsp.Extension"

	// everything else is base.String() with
	// the first letter as lowercase
	default:
		return strings.ToLower(s.BaseName())
	}
}
func (k Primitive) String() string {
	switch k {
	case String:
		return "String"
	case Bytes:
		return "Bytes"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Complex64:
		return "Complex64"
	case Complex128:
		return "Complex128"
	case Uint:
		return "Uint"
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Byte:
		return "Byte"
	case Int:
		return "Int"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Bool:
		return "Bool"
	case Intf:
		return "Intf"
	case Time:
		return "time.Time"
	case Ext:
		return "Extension"
	case IDENT:
		return "Ident"
	default:
		return "INVALID"
	}
}

```