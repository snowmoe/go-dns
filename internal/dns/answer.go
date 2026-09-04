package dns

type Resource struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   []byte
}

func (c *Cursor) DecodeResource() Resource {
	var (
		name   = c.parseName()
		rType  = c.parseType()
		class  = c.parseClass()
		ttl    = c.parseTTL()
		length = c.parseRDLength()
		data   = c.parseRData(length)
	)

	return Resource{
		Name:   name,
		Type:   rType,
		Class:  class,
		TTL:    ttl,
		Length: length,
		Data:   data,
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

func (c *Cursor) parseType() uint16     { return c.takeUint16() }
func (c *Cursor) parseClass() uint16    { return c.takeUint16() }
func (c *Cursor) parseTTL() uint32      { return c.takeUint32() }
func (c *Cursor) parseRDLength() uint16 { return c.takeUint16() }

func (c *Cursor) parseRData(l uint16) []byte {
	b := c.takeBytes(int(l))
	c.Pos += int(l)
	return b
}
