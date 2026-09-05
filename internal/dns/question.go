package dns

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

func (c *Cursor) DecodeQuestion() Question {
	return Question{
		Name:  c.parseLabel(),
		Type:  c.takeUint16(),
		Class: c.takeUint16(),
	}
}
