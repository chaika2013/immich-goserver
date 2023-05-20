package helper

import (
	"hash"
	"hash/crc32"
	"io"
)

func NewCRC32Writer(poly uint32, w io.Writer) *CRC32Writer {
	return &CRC32Writer{
		h: crc32.New(crc32.MakeTable(poly)),
		w: w,
	}
}

type CRC32Writer struct {
	h hash.Hash32
	w io.Writer
}

func (c *CRC32Writer) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p) // with each write ...
	c.h.Write(p)          // ... update the hash
	return
}

// Sum gives the final hash
func (c *CRC32Writer) Sum() uint32 {
	return c.h.Sum32()
}
