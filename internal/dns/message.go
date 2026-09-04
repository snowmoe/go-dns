package dns

type Message struct {
	Header    Header
	Questions []Question
	Resources []Resource
}

func (c *Cursor) DecodeMessage() Message {
	h := c.DecodeHeader()

	qCount := int(h.QDCount)
	rCount := int(h.ANCount)

	q := make([]Question, qCount)
	r := make([]Resource, rCount)

	for i := range q {
		q[i] = c.DecodeQuestion()
	}

	for i := range r {
		r[i] = c.DecodeResource()
	}

	return Message{
		Header:    h,
		Questions: q,
		Resources: r,
	}
}
