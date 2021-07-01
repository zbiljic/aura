package ulidx

import (
	cryptorand "crypto/rand"
	mathrand "math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func ULID() ulid.ULID {
	entropy := cryptorand.Reader

	return ulid.MustNew(ulid.Now(), entropy)
}

//nolint:gosec
func QuickULID() ulid.ULID {
	seed := time.Now().UnixNano()
	source := mathrand.NewSource(seed)
	entropy := mathrand.New(source)

	return ulid.MustNew(ulid.Now(), entropy)
}

func ZeroULID() ulid.ULID {
	entropy := zeroReader{}

	return ulid.MustNew(ulid.Now(), entropy)
}

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}

	return len(p), nil
}
