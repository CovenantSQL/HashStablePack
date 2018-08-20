package covenant

import (
	"bytes"
	"testing"
	"time"
	"gitlab.com/thunderdb/ThunderDB/crypto/hash"
)

// test different type and member name but same data type and content hash identical
func TestMarshalHashAccountStable2(t *testing.T) {
	tm := time.Now()
	v1 := HeaderTest{
		Version:     110,
		TestName:    "31231",
		TestArray:   []byte{0x11, 0x22},
		Producer:    "rewqrwe",
		GenesisHash: []hash.Hash{
			{0x10},
			{0x20},
		},
		ParentHash:  []*hash.Hash{
			{0x10},
			{0x20},
		},
		MerkleRoot:  &[]*hash.Hash{
			{0x10},
			{0x20},
		},
		Timestamp:   tm,
	}
	v2 := HeaderTest2{
		Version2:     110,
		TestName2:    "31231",
		TestArray2:   []byte{0x11, 0x22},
		Producer2:    "rewqrwe",
		GenesisHash2: []hash.Hash{
			{0x10},
			{0x20},
		},
		ParentHash2:  []*hash.Hash{
			{0x10},
			{0x20},
		},
		MerkleRoot2:  &[]*hash.Hash{
			{0x10},
			{0x20},
		},
		Timestamp2:   tm,
	}
	bts1, err := v1.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	bts2, err := v2.MarshalHash()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bts1, bts2) {
		t.Fatal("hash not stable")
	}
}

// test different type and member name but same data type and content hash identical
func TestMarshalHashAccountStable3(t *testing.T) {
	p1 := Person1{
		Name1:       "Auxten",
		Address1:    "@CovenantSQL.io",
		Age1:        70,
		unexported1: false,
	}
	p2 := Person2{
		Name2:       "Auxten",
		Address2:    "@CovenantSQL.io",
		Age2:        70,
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
