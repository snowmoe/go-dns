package dns

type Cursor struct {
	Pos    int
	Data   []byte
	Labels map[int]string
}
