package helpers

import (
	"io"

	"github.com/gotd/td/bin"
)

func RandInt64(randSource io.Reader) (int64, error) {
	var buf [bin.Word * 2]byte
	if _, err := io.ReadFull(randSource, buf[:]); err != nil {
		return 0, err
	}
	b := &bin.Buffer{Buf: buf[:]}
	return b.Long()
}
