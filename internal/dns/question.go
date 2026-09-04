package dns

import (
	"strings"
)

func (c *Cursor) DecodeQuestion() Question {
	return Question{
		Name:  c.parseQName(),
		Type:  c.parseQType(),
		Class: c.parseQClass(),
	}
}

func (c *Cursor) parseQType() uint16  { return c.takeUint16() }
func (c *Cursor) parseQClass() uint16 { return c.takeUint16() }

func (c *Cursor) parseQName() string {
	var (
		parts []string
		start = c.Pos
	)

	for {
		len := int(c.Data[c.Pos])

		if len == 0 {
			c.Pos += 1
			break
		}

		parts = append(parts, string(c.Data[c.Pos+1:c.Pos+len+1]))
		c.Pos = c.Pos + len + 1
	}

	label := strings.Join(parts, ".")
	c.Labels[start] = label

	return label
}
