// Package hash stands in for CovenantSQL's crypto/hash in the generator
// tests, so that this module does not require CovenantSQL.
package hash

import hsp "github.com/CovenantSQL/HashStablePack/v2/marshalhash"

// Hash is a fixed-size hash with a hand-written MarshalHash, encoded as
// CovenantSQL's crypto/hash.Hash is.
type Hash [32]byte

// MarshalHash returns a copy of the raw hash bytes.
func (h *Hash) MarshalHash() (o []byte, err error) {
	return append([]byte(nil), h[:]...), nil
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message.
func (h *Hash) Msgsize() (s int) {
	return hsp.BytesPrefixSize + len(h)
}
