package dns

import (
	"encoding/binary"
	"fmt"
)

func (c *Cursor) takeBytes(n int) []byte {
	b := c.Data[c.Pos : c.Pos+n]
	c.Pos += n
	return b
}

func (c *Cursor) takeUint16() uint16 {
	b := c.takeBytes(2)
	return binary.BigEndian.Uint16(b)
}

func (c *Cursor) takeUint32() uint32 {
	b := c.takeBytes(4)
	return binary.BigEndian.Uint32(b)
}

func bits(value uint16, shift, width int) uint16 {
	return (value >> shift) & ((1 << width) - 1)
}

func (c *Cursor) skipBytes(n int) {
	c.Pos += n
}

func (c *Cursor) PrintPos() {
	fmt.Println("position:", c.Pos)
}
