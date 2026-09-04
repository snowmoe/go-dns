package dns

import "encoding/binary"

func (c *Cursor) ReadResource() Resource {
	name := c.parseName()

	c.SkipBytes(8)

	rdl := c.parseRDLength()
	rd := c.parseRData(rdl)

	return Resource{
		Name:   name,
		Length: rdl,
		Data:   rd,
	}
}

func (c *Cursor) parseName() string {
	b := c.Data[c.Pos]

	// top two bytes = 11 = pointer
	if b&0xC0 == 0xC0 {
		label := c.Labels[int(c.Data[c.Pos+1])]
		c.Pos += 2
		return label
	} else {
		return ""
	}
}

func (c *Cursor) parseRDLength() uint16 {
	b := c.takeBytes(2)
	return binary.BigEndian.Uint16(b)
}

func (c *Cursor) parseRData(l uint16) []byte {
	b := c.takeBytes(int(l))
	c.Pos += int(l)
	return b
}
