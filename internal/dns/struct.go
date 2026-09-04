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

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

type Resource struct {
	Name   string
	Type   uint16
	Class  uint16
	TTL    uint32
	Length uint16
	Data   []byte
}
