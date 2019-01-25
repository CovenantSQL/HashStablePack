##  HSP 为什么比 golang 的 reflect.DeepEqual 快10倍



使用 Golang 的时候遇到一个棘手的问题：

假设我们现在有一个比较复杂的结构体 ComplexStruct

```go
type SimpleStruct struct {
    Str string
}
type ComplexStruct struct {
    Simple map[string]*SimpleStruct
    Other []byte   
    Nums  [8]float64 
}
```



我们需要在 10000 个 ComplexStruct 类型的变量 Slice 里找到和已知变量内容相同的。

首先，最笨的办法当然是手写:

```go
var target = ComplexStruct {
    ...
}
for _, cs := range csSlice {
    for k, v := range cs.Simple {
        if v1, ok := target.Simple[k]; ok {
             // 算了老子不干了
        } else {
            
        }
    }
}
```



## DeepEqual

聪明的朋友都会用 `reflect.DeepEqual` 来比较：

```go
reflect.DeepEqual(v1, v2)
```

但我们知道，无论是哪种语言，基于反射（reflect）的各种方法都会比普通的函数调用至少慢上 1 ~ 2 个数量级。

有没有又懒又快的方法呢？我反正是没找到，所以自己基于[一个 MsgPack 的库](https://github.com/tinylib/msgp)写了一个 [HashStatePack](https://github.com/CovenantSQL/HashStablePack)，主要功能如下：

1. 根据已经定义好的 golang 结构体自动生成代码，避免搬砖

2. 比`reflect.DeepEqual`快大概10 ~ 20 倍，[测试在此](https://github.com/CovenantSQL/HashStablePack/blob/master/test/hashstable_test.go#L106)

   ```go
   BenchmarkCompare/benchmark_reflect-8          20074 ns/op //reflect.DeepEqual
   BenchmarkCompare/benchmark_hsp-8               2322 ns/op
   BenchmarkCompare/benchmark_hsp_1_cached-8      1101 ns/op
   BenchmarkCompare/benchmark_hsp_both_cached-8   11.2 ns/op
   ```

## 为什么不用……

为什么不用 Protobuf 或者 MsgPack，甚至 JSON 序列化之后再比较？

- JSON: 内存使用效率太低，特别是遇到 Binary 类型。
- 大部分 JSON、MsgPack 的库也是基于反射的实现，也很慢。
- Prorobuf: struct must defined in proto language, and other limitations discussed [here](https://gist.github.com/kchristidis/39c8b310fd9da43d515c4394c3cd9510)
- 最后，也是最棘手的问题：Golang 的设计者为了避免大家错误的依赖 map 的顺序，在迭代 map 的时候故意加入了一定的洗牌算法。这就导致几乎针对同样一个 map 的 range 每次的结果都不一样。

### 原理

1. 读取 .go 源代码文件，生成 AST（抽象语法树）
2. 找到所有的类型，排序
3. 对每个 Struct 内部成员按照 tag （`hsp:"xxx"`）进行排序，没有 tag 的用名字
4. 根据不同的类型生成不同的序列化代码
5. 生成测试代码

例如，我们开头例子中的 ComplexStruct，生成的核心代码如下：

```go
// MarshalHash marshals for hash
func (z *ComplexStruct) MarshalHash() (o []byte, err error) {
   var b []byte
   o = hsp.Require(b, z.Msgsize())
   // map header, size 3
   o = append(o, 0x83)
   o = hsp.AppendArrayHeader(o, uint32(8))
   for za0003 := range z.Nums {
      o = hsp.AppendFloat64(o, z.Nums[za0003])
   }
   o = hsp.AppendBytes(o, z.Other)
   o = hsp.AppendMapHeader(o, uint32(len(z.Simple)))
   za0001Slice := make([]string, 0, len(z.Simple))
   for i := range z.Simple {
      za0001Slice = append(za0001Slice, i)
   }
   sort.Strings(za0001Slice)
   for _, za0001 := range za0001Slice {
      za0002 := z.Simple[za0001]
      o = hsp.AppendString(o, za0001)
      if za0002 == nil {
         o = hsp.AppendNil(o)
      } else {
         // map header, size 1
         o = append(o, 0x81)
         o = hsp.AppendString(o, za0002.Str)
      }
   }
   return
}
```

剩下的代码就可以这么写了：

```
bts1, _ := v1.MarshalHash()
bts2, _ := v2.MarshalHash()
if bytes.Equal(bts1, bts2) {
    ...
}
```

针对我们遇到的在大量 Slice 中 寻找相同内容的问题，如果我们对生成的`[]byte`，进行一次哈希。然后用哈希只作为 key，对象作为 Value，效率将会非常的高。

HashStablePack 目前主要被 [CovenantSQL](https://github.com/CovenantSQL/CovenantSQL) 用来做签名、校验，以及区块哈希计算上，希望可以帮到你:-)



## 怎么使用

```go
go get -u github.com/CovenantSQL/HashStablePack/hsp
```

在你需要生成的源文件头部加上

```go
//go:generate hsp
```

运行

```go
go generate ./...
```

代码、测试代码，就统统生成好了