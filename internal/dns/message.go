package dns

func (c *Cursor) DecodeMessage() Message {
	return Message{
		Header:   c.DecodeHeader(),
		Question: c.DecodeQuestion(),
		Resource: c.DecodeResource(),
	}
}
