package dns

type Cursor struct {
	Pos    int
	Data   []byte
	Labels map[int]string
}

type Message struct {
	Header   Header
	Question Question
	Resource Resource
}
