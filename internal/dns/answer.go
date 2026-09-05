package dns

import "strings"

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
	if c.isPointer() {
		// we know it's a pointer so we skip that byte
		c.skipBytes(1)

		// consume the next byte and make note of our position
		var (
			pointer = c.takeByte()
			current = c.pos()
		)

		// go back to the offset from the pointer byte
		c.Pos = int(pointer)

		// parse as normal
		label := c.parseLabel()

		// go forward to where we were
		c.Pos = current

		return strings.Join(label, ".")
	} else {
		label := c.parseLabel()

		return strings.Join(label, ".")
	}
}

func (c *Cursor) parseLabel() []string {
	var parts []string
	for {
		labelLen := c.takeByte()

		if labelLen == 0 {
			break
		}

		label := string(c.takeBytes(int(labelLen)))

		parts = append(parts, label)
	}

	return parts
}

func (c *Cursor) parseRData(l uint16) []byte {
	b := c.takeBytes(int(l))
	return b
}
