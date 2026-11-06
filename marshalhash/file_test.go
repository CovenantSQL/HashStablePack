// +build linux darwin dragonfly freebsd netbsd openbsd

package marshalhash_test

import (
	"bytes"
	"crypto/rand"
	"github.com/CovenantSQL/HashStablePack/v2/marshalhash"
	"os"
	"testing"
)

type rawBytes []byte

func (r rawBytes) MarshalMsg(b []byte) ([]byte, error) {
	return marshalhash.AppendBytes(b, []byte(r)), nil
}

func (r rawBytes) Msgsize() int {
	return marshalhash.BytesPrefixSize + len(r)
}

func (r *rawBytes) UnmarshalMsg(b []byte) ([]byte, error) {
	tmp, out, err := marshalhash.ReadBytesBytes(b, (*(*[]byte)(r))[:0])
	*r = rawBytes(tmp)
	return out, err
}

func TestReadWriteFile(t *testing.T) {
	t.Parallel()

	f, err := os.Create("tmpfile")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		f.Close()
		os.Remove("tmpfile")
	}()

	data := make([]byte, 1024*1024)
	rand.Read(data)

	err = marshalhash.WriteFile(rawBytes(data), f)
	if err != nil {
		t.Fatal(err)
	}

	var out rawBytes
	f.Seek(0, os.SEEK_SET)
	err = marshalhash.ReadFile(&out, f)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte(out), []byte(data)) {
		t.Fatal("Input and output not equal.")
	}
}
