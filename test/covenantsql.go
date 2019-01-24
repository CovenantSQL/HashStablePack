package covenant

import (
	"time"

	"github.com/CovenantSQL/CovenantSQL/crypto/hash"
	"github.com/CovenantSQL/CovenantSQL/proto"
)

//go:generate hsp

const Eight = 8

type MyInt int
type Data []byte

type Struct struct {
	Which map[string]*MyInt `hsp:"2which"`
	Other Data              `hsp:"1other"`
	Nums  [Eight]float64    `hsp:"3nums"`
}

// HeaderTest is a block header.
type HeaderTest struct {
	Version     int32         `hsp:"01"`
	TestName    string        `hsp:"00"`
	TestArray   []byte
	Producer    proto.NodeID  `hsp:"02"`
	GenesisHash []hash.Hash   `hsp:"06"`
	ParentHash  []*hash.Hash  `hsp:"03"`
	MerkleRoot  *[]*hash.Hash `hsp:"05"`
	Timestamp   time.Time     `hsp:"04"`
	xx int
}

// HeaderTest is a block header.
type HeaderTest2 struct {
	Version2     int32         `hsp:"01"`
	TestName2    string        `hsp:"00"`
	TestArray2   []byte
	Producer2    proto.NodeID  `hsp:"02"`
	GenesisHash2 []hash.Hash   `hsp:"06"`
	ParentHash2  []*hash.Hash  `hsp:"03"`
	MerkleRoot2  *[]*hash.Hash `hsp:"05"`
	Timestamp2   time.Time     `hsp:"04"`
	xx int
}

type Person1 struct {
	Name1       string
	Age1        int
	Address1    string
	Map         map[string]int
	unexported1 bool // this field is ignored
}

type Person2 struct {
	Name2       string
	Address2    string
	Age2        int
	Map         map[string]int
	unexported2 bool // this field is ignored
}
