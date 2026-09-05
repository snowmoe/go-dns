package dns

import (
	"strings"
)

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (c *Cursor) DecodeQuestion() Question {
	return Question{
		Name:  c.parseQName(),
		Type:  c.takeUint16(),
		Class: c.takeUint16(),
	}
}

func (c *Cursor) parseQName() string {
	if c.isPointer() {
		c.skipBytes(1)
		pointer := c.takeByte()
		current := c.pos()

		c.Pos = int(pointer)

		label := c.parseLabel()

		c.Pos = current

		return strings.Join(label, ".")
	} else {
		label := c.parseLabel()

		return strings.Join(label, ".")
	}
}
