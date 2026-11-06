// Package proto stands in for CovenantSQL's proto in the generator tests,
// so that this module does not require CovenantSQL.
package proto

import hsp "github.com/CovenantSQL/HashStablePack/v2/marshalhash"

// NodeID is a string type encoded as CovenantSQL's proto.NodeID is.
type NodeID string

// MarshalHash marshals for hash.
func (z NodeID) MarshalHash() (o []byte, err error) {
	return hsp.AppendString(hsp.Require(nil, z.Msgsize()), string(z)), nil
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message.
func (z NodeID) Msgsize() (s int) {
	return hsp.StringPrefixSize + len(z)
}
