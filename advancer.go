package midec

import (
	"io"
)

const tmpLength = 256 * 3

// ReadAdvancer is the struct that can skip some bytes reading.
type ReadAdvancer struct {
	io.Reader
	tmp []byte
}

// NewReadAdvancer creates ReadAdvancer.
func NewReadAdvancer(r io.Reader) *ReadAdvancer {
	var arr [tmpLength]byte
	return &ReadAdvancer{
		Reader: r,
		tmp:    arr[:],
	}
}

// ReadFull is a shorthand for io.ReadFull.
func (a *ReadAdvancer) ReadFull(buf []byte) (int, error) {
	return io.ReadFull(a.Reader, buf)
}

// Advance skips some bytes.
func (a *ReadAdvancer) Advance(n uint) error {
	for n >= tmpLength {
		buf := a.tmp[0:tmpLength]
		_, err := a.ReadFull(buf)
		if err != nil {
			return err
		}

		n -= tmpLength
	}

	if n > 0 {
		buf := a.tmp[0:n]
		_, err := a.ReadFull(buf)
		if err != nil {
			return err
		}
	}
	return nil
}
