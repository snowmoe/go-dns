package dns

import "fmt"

func (c *Cursor) takeBytes(n int) []byte {
	b := c.Data[c.Pos : c.Pos+n]
	c.Pos += n
	return b
}

func (c *Cursor) SkipBytes(n int) {
	c.Pos += n
}

func (c *Cursor) PrintPos() {
	fmt.Println("position:", c.Pos)
}
