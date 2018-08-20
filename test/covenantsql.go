package covenant

import (
	"gitlab.com/thunderdb/ThunderDB/proto"
	"gitlab.com/thunderdb/ThunderDB/crypto/hash"
	"time"
)

//go:generate hsp

const Eight = 8
type MyInt int
type Data []byte

type Struct struct {
	Which  map[string]*MyInt `msg:"which"`
	Other  Data              `msg:"other"`
	Nums   [Eight]float64    `msg:"nums"`
}

// HeaderTest is a block header.
type HeaderTest struct {
	Version     int32
	TestName    string
	TestArray   []byte
	Producer    proto.NodeID
	GenesisHash []hash.Hash
	ParentHash  []*hash.Hash
	MerkleRoot  *[]*hash.Hash
	Timestamp   time.Time
}

// HeaderTest is a block header.
type HeaderTest2 struct {
	Version2     int32
	TestName2    string
	TestArray2   []byte
	Producer2    proto.NodeID
	GenesisHash2 []hash.Hash
	ParentHash2  []*hash.Hash
	MerkleRoot2  *[]*hash.Hash
	Timestamp2   time.Time
}

type Person1 struct {
	Name1       string
	Address1    string
	Age1        int
	unexported1 bool             // this field is ignored
}

type Person2 struct {
	Name2       string
	Address2    string
	Age2        int
	unexported2 bool             // this field is ignored
}
