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
		name   = c.parseLabel()
		rType  = c.takeUint16()
		class  = c.takeUint16()
		ttl    = c.takeUint32()
		length = c.takeUint16()
		data   = c.takeBytes(int(length))
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

func (c *Cursor) parseLabel() string {
	var parts []string
	for {
		// check if the current byte is a pointer
		if c.isPointer() {
			// parse the label from the pointer
			label := c.parsePointerLabel()
			parts = append(parts, label)
			// break because pointers can only be at the end of label section
			break
		}

		// wasn't a pointer so we know it's the legnth of label bytes
		labelLen := c.takeByte()

		// check if the length = 0 = terminator
		if labelLen == 0 {
			break
		}

		label := string(c.takeBytes(int(labelLen)))

		parts = append(parts, label)
	}

	return strings.Join(parts, ".")
}

func (c *Cursor) parsePointerLabel() string {
	// pointer offset spreads over 14 bits
	pointer := bits(c.takeUint16(), 0, 14)
	current := c.pos()

	// go back to the offset from the pointer byte
	c.Pos = int(pointer)

	// parse as normal
	label := c.parseLabel()

	// go forward to where we were
	c.Pos = current

	return label
}
