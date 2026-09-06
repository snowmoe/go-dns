package dns

import (
	"net"
	"strings"
)

type Resource struct {
	Name   string
	Type   Type
	Class  Class
	TTL    uint32
	Length uint16
	Data   any
}

func (c *Cursor) DecodeResource() Resource {
	var (
		name   = c.parseLabel()
		rType  = c.takeUint16()
		class  = c.takeUint16()
		ttl    = c.takeUint32()
		length = c.takeUint16()
		data   = c.parseRData(rType, length)
	)

	return Resource{
		Name:   name,
		Type:   Type(rType),
		Class:  Class(class),
		TTL:    ttl,
		Length: length,
		Data:   data,
	}
}

func (c *Cursor) parseRData(t uint16, ul uint16) any {
	l := int(ul)

	switch t {
	case 1:
		return net.IP(c.takeBytes(l))
	case 2:
		return c.parseLabel()
	default:
		return make([]byte, l)
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

func (c Class) String() string {
	switch c {
	case 1:
		return "IN"
	case 2:
		return "CS"
	case 3:
		return "CH"
	case 4:
		return "HS"
	case 255:
		return "*"
	default:
		return ""
	}
}

func (t Type) String() string {
	switch t {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 5:
		return "CNAME"
	case 12:
		return "PTR"
	case 15:
		return "MX"
	case 16:
		return "TXT"
	default:
		return ""
	}
}
