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
		rType  = c.takeUint16()
		class  = c.takeUint16()
		ttl    = c.takeUint32()
		length = c.takeUint16()
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

	// is a pointer
	if b&0xC0 == 0xC0 {
		// offset := c.Data[c.Pos+1]
		c.Pos += 2
		return "POINTER"

	} else {
		for {
			len := int(c.Data[c.Pos])
			if len == 0 {
				c.Pos += 1
				break
			}

			c.skipBytes(len + 1)
		}

		return "NOTPOINTER"
	}
}

func (c *Cursor) parseRData(l uint16) []byte {
	b := c.takeBytes(int(l))
	return b
}
