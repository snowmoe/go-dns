package dns

type Question struct {
	Name  string
	Type  Type
	Class Class
}

func (c *Cursor) DecodeQuestion() Question {
	return Question{
		Name:  c.parseLabel(),
		Type:  Type(c.takeUint16()),
		Class: Class(c.takeUint16()),
	}
}
