package covenant

import (
	"gitlab.com/thunderdb/ThunderDB/proto"
	"gitlab.com/thunderdb/ThunderDB/utils"
	"gitlab.com/thunderdb/ThunderDB/crypto/hash"
	"time"
	"bytes"
	"encoding/binary"
)

//go:generate hsp

// HeaderTest is a block header.
type HeaderTest struct {
	Version     int32
	TestName	string
	TestArray 	[]byte
	Producer    proto.NodeID
	GenesisHash []hash.Hash
	ParentHash  []*hash.Hash
	MerkleRoot  *[]*hash.Hash
	Timestamp   time.Time
}

// MarshalHash marshals for hash
func (h *HeaderTest) MarshalHash() ([]byte, error) {
	buffer := bytes.NewBuffer(nil)

	if err := utils.WriteElements(buffer, binary.BigEndian,
		h.Version,
		h.Producer,
		&h.GenesisHash,
		&h.ParentHash,
		&h.MerkleRoot,
		h.Timestamp,
	); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

